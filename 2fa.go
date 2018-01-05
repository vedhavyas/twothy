package twofa

import (
	"encoding/binary"
)

const (
	// DefaultStepTime as per RFC 6238 is 30 seconds
	DefaultStepTime int = 30

	// DefaultT0 as per RFC 6238 is 0
	DefaultT0 int = 0

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
