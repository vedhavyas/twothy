package twothy

import (
	"bytes"
	"testing"
)

func Test_addPadding(t *testing.T) {
	tests := []struct {
		a []byte
		r int
		b []byte
	}{
		{
			a: nil,
			r: 1,
			b: []byte{passwordPad},
		},

		{
			a: []byte{100, 255, 200},
			r: 3,
			b: []byte{100, 255, 200, passwordPad, passwordPad, passwordPad},
		},

		{
			a: nil,
			r: 0,
			b: nil,
		},

		{
			a: []byte{},
			r: 0,
			b: []byte{},
		},
	}

	for _, c := range tests {
		r := addPadding(c.a, c.r)
		if !bytes.Equal(c.b, r) {
			t.Fatalf("expected %v got %v", c.b, r)
		}
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		pwd    int
		length int
		err    bool
	}{
		{
			pwd:    0,
			length: 16,
		},

		{
			pwd:    15,
			length: 16,
		},

		{
			pwd:    18,
			length: 24,
		},

		{
			pwd:    27,
			length: 32,
		},

		{
			pwd: 34,
			err: true,
		},
	}

	for _, c := range tests {
		var password []byte = nil
		for i := 0; i < c.pwd; i++ {
			password = append(password, 'i')
		}

		r, err := validatePassword(password)
		if err != nil {
			if c.err {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if c.length != len(r) {
			t.Fatalf("expected password length %d but got %d", c.length, len(r))
		}
	}
}

func Test_encrypt_decrypt(t *testing.T) {
	tests := []string{
		"hello, world!!",
		"This is a long messages",
		"",
		"take this!!!!",
		"this is a even longer message. Longer longer longer!!!!",
	}

	password := addPadding([]byte("password"), 8)

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
