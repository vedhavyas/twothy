package twothy

import (
	"testing"
)

func Test_encrypt_decrypt(t *testing.T) {
	tests := []string{
		"hello, world!!",
		"This is a long messages",
		"",
		"take this!!!!",
		"this is a even longer message. Longer longer longer!!!!",
	}

	password := hashOf([]byte("password"))

	for _, c := range tests {
		ed, err := encrypt([]byte(c), password)
		if err != nil {
			t.Fatalf("unexpected error while encrypting: %v", err)
		}

		dd, err := decrypt(ed, password)
		if err != nil {
			t.Fatalf("unexpected error while decrypting: %v", err)
		}

		if string(dd) != c {
			t.Fatalf("expected '%s' but got '%s'", c, string(dd))
		}
	}
}
