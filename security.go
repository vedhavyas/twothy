package twothy

import (
	"fmt"

	"golang.org/x/crypto/ssh/terminal"
)

// passwordPad to add extra bytes
const passwordPad = byte('=')

// getPassword returns the password taken from the user
func getPassword() ([]byte, error) {
	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read user's password: %v", err)
	}

	return validatePassword(password)
}

// addPadding repeat adds the passwordPad to pwd
func addPadding(pwd []byte, repeat int) []byte {
	for i := 0; i < repeat; i++ {
		pwd = append(pwd, passwordPad)
	}

	return pwd
}

// validatePassword ensures password is either 16, 24, 32 bytes long
func validatePassword(pwd []byte) ([]byte, error) {
	l := len(pwd)
	switch {
	case l <= 16:
		return addPadding(pwd, 16-l), nil
	case l <= 24:
		return addPadding(pwd, 24-l), nil
	case l <= 32:
		return addPadding(pwd, 32-l), nil
	default:
		return nil, fmt.Errorf("password should be <=32 characters: %d", len(pwd))
	}
}
