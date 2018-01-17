package twothy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/ssh/terminal"
)

// hashOf will return a 256-bit hash of pwd
func hashOf(pwd []byte) []byte {
	h := sha256.Sum256(pwd)
	return h[:]
}

// getPassword returns the password taken from the user
func getPassword(hint string) ([]byte, error) {
	fmt.Printf("Password(%s): ", hint)
	password, err := terminal.ReadPassword(0)
	if err != nil {
		return nil, fmt.Errorf("failed to read user's password: %v", err)
	}

	if len(password) < 1 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	return hashOf(password), nil
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
