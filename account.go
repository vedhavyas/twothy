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
// Issuer and Label together provides a unique account even for same service.
// Ex: google with multiple accounts
type Account struct {
	Issuer   string `json:"issuer"`
	Label    string `json:"label"`
	T0       int    `json:"t_0"`
	StepTime int    `json:"step_time"`
	Digits   int    `json:"digits"`
	Key      string `json:"key"`
}

// NewAccount returns a new 2fa account with default values
func NewAccount(issuer, label, key string) Account {
	key = strings.Replace(key, "-", "", -1)
	key = strings.Replace(key, " ", "", -1)
	key = strings.ToUpper(key)
	return Account{
		Issuer:   issuer,
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
func CreateOTP(a Account, time int64) (otp string, err error) {
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

	otp = fmt.Sprint(r % int32(pow(10, a.Digits)))
	if len(otp) != a.Digits {
		rpt := strings.Repeat("0", a.Digits-len(otp))
		otp = rpt + otp
	}

	return otp, nil
}

// SaveAccount writes account info to twothy folder
func SaveAccount(c Config, a Account, pwd []byte) error {
	fileName := fmt.Sprintf("%s_%s.twothy", a.Issuer, a.Label)
	path := fmt.Sprintf("%s%s", c.AccountsFolder, fileName)
	data, err := json.Marshal(a)
	if err != nil {
		return fmt.Errorf("failed to marshall the object: %v", err)
	}

	ed, err := encrypt(data, pwd)
	if err != nil {
		return fmt.Errorf("failed to encrypt account details: %v", err)
	}

	err = writeToFile(path, []byte(ed))
	if err != nil {
		return fmt.Errorf("failed to save account: %v", err)
	}

	return nil
}

// loadAccount will load the account from file
func loadAccount(filePath string, pwd []byte) (a Account, err error) {
	ed, err := ioutil.ReadFile(filePath)
	if err != nil {
		return a, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	data, err := decrypt(string(ed), pwd)
	if err != nil {
		return a, fmt.Errorf("failed to decrypt the account %s: %v", filePath, err)
	}

	err = json.Unmarshal(data, &a)
	if err != nil {
		return a, fmt.Errorf("failed to load file %s: %v", filePath, err)
	}

	return a, nil
}

// LoadAccounts will a load accounts matching issuer and label
// if label is empty, loads all the accounts matching issuer
// if issuer and label are empty, all th accounts are returned
func LoadAccounts(c Config, issuer, label string, pwd []byte) (accounts []Account, err error) {
	if issuer != "" && label != "" {
		path := fmt.Sprintf("%s%s_%s.twothy", c.AccountsFolder, issuer, label)
		a, err := loadAccount(path, pwd)
		return []Account{a}, err
	}

	err = filepath.Walk(c.AccountsFolder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || filepath.Ext(path) != ".twothy" {
			return nil
		}

		a, err := loadAccount(path, pwd)
		if err != nil {
			return fmt.Errorf("failed to read account in %s: %v", path, err)
		}

		if issuer != "" && strings.ToLower(issuer) != strings.ToLower(a.Issuer) {
			return nil
		}

		accounts = append(accounts, a)
		return nil
	})

	if err != nil {
		return accounts, fmt.Errorf("failed to load accounts: %v", err)
	}

	return accounts, nil
}
