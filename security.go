package twothy

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

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
	pad := []byte{passwordPad}
	pad = bytes.Repeat(pad, repeat)
	return append(pwd, pad...)
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

// encrypt encrypts the message with AES
func encrypt(msg, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		// this should not happen
		return "", err
	}

	ct := make([]byte, aes.BlockSize+len(msg))
	iv := ct[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ct[aes.BlockSize:], msg)
	return base64.URLEncoding.EncodeToString(ct), nil
}

// decrypt decrypts the cipher text to plaintext
func decrypt(msg string, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decodedMsg, err := base64.URLEncoding.DecodeString(msg)
	if err != nil {
		return nil, err
	}

	iv := decodedMsg[:aes.BlockSize]
	data := decodedMsg[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(data, data)

	return data, nil
}
