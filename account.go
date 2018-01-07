package twothy

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultStepTime as per RFC 6238 is 30 seconds
	DefaultStepTime = 30

	// DefaultT0 as per RFC 6238 is 0
	DefaultT0 = 0

	// DefaultDigits as per RFC 6238 is 6 digits
	DefaultDigits = 6
)

// Account represents a single 2fa activated account
// Name and Label together provides a unique account even for same service.
// Ex: google with multiple accounts
type Account struct {
	Name     string
	Label    string
	T0       int
	StepTime int
	Digits   int
	Key      string
}

// NewAccount returns a new 2fa account with default values
func NewAccount(name, label, key string) Account {
	return Account{
		Name:     name,
		Label:    label,
		T0:       DefaultT0,
		StepTime: DefaultStepTime,
		Digits:   DefaultDigits,
		Key:      key,
	}
}

// getMessage constructs the message for HMAC with given params
func getMessage(t1 int64, t0, stepTime int) (message []byte) {
	step := (t1 - int64(t0)) / int64(stepTime)
	message = make([]byte, 8)
	binary.BigEndian.PutUint64(message, uint64(step))
	return message
}

// getOffset returns the offset from hash bytes as per https://tools.ietf.org/html/rfc4226#section-5.4
func getOffset(hash []byte) int {
	lastByte := hash[len(hash)-1]
	return int(lastByte & 0xf)
}

// pow returns x^y
func pow(x, y int) int {
	if x == 0 || x == 1 {
		return x
	}

	if y == 0 {
		return 1
	}

	r := 1
	for i := 0; i < y; i++ {
		r *= x
	}

	return r
}

// CreateOTP creates opt for account at given time
func CreateOTP(a Account, time int64) (otp int32, err error) {
	m := getMessage(time, a.T0, a.StepTime)

	key, err := base32.StdEncoding.DecodeString(a.Key)
	if err != nil {
		return otp, fmt.Errorf("invalid key: %v", err)
	}

	hasher := hmac.New(sha1.New, key)
	_, err = hasher.Write(m)
	if err != nil {
		return otp, fmt.Errorf("failed create hash: %v", err)
	}

	h := hasher.Sum(nil)
	ofs := getOffset(h)
	r := int32(h[ofs]&0x7f)<<24 |
		int32(h[ofs+1])<<16 |
		int32(h[ofs+2])<<8 |
		int32(h[ofs+3])

	return r % int32(pow(10, a.Digits)), nil
}

// SaveAccount writes account info to twothy folder
func SaveAccount(c Config, a Account) error {
	fileName := fmt.Sprintf("%s_%s.twothy", a.Name, a.Label)
	path := fmt.Sprintf("%s%s", c.AccountsFolder, fileName)
	err := writeToFile(path, a)
	if err != nil {
		return fmt.Errorf("failed to save account: %v", err)
	}

	return nil
}

// loadAccount will load the account from file
func loadAccount(filePath string) (a Account, err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return a, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &a)
	if err != nil {
		return a, fmt.Errorf("failed to load file %s: %v", filePath, err)
	}

	return a, nil
}

// LoadAccounts will a load accounts matching name and label
// if label is empty, loads all the accounts matching name
// if name and label are empty, all th accounts are returned
func LoadAccounts(c Config, name, label string) (accounts []Account, err error) {
	if name != "" && label != "" {
		path := fmt.Sprintf("%s%s_%s.twothy", c.AccountsFolder, name, label)
		a, err := loadAccount(path)
		return []Account{a}, err
	}

	err = filepath.Walk(c.AccountsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".twothy" {
			a, err := loadAccount(path)
			if err != nil {
				return fmt.Errorf("failed to read account in %s: %v", path, err)
			}

			if name != "" && strings.ToLower(name) != strings.ToLower(a.Name) {
				return nil
			}

			accounts = append(accounts, a)
		}
		return nil
	})

	if err != nil {
		return accounts, fmt.Errorf("failed to load accounts: %v", err)
	}

	return accounts, nil
}
