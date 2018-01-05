package twofa

import (
	"bytes"
	"testing"
)

func Test_getMessage(t *testing.T) {
	tests := []struct {
		t0       int
		t1       int64
		stepTime int
		result   []byte
	}{
		{
			t0:       0,
			t1:       1257894000,
			stepTime: 30,
			result:   []byte{0x0, 0x0, 0x0, 0x0, 0x2, 0x7f, 0xcc, 0x48},
		},

		{
			t0:       0,
			t1:       1257894020,
			stepTime: 30,
			result:   []byte{0x0, 0x0, 0x0, 0x0, 0x2, 0x7f, 0xcc, 0x48},
		},

		{
			t0:       0,
			t1:       1515150425,
			stepTime: 30,
			result:   []byte{0x0, 0x0, 0x0, 0x0, 0x3, 0x2, 0xa5, 0x36},
		},

		{
			t0:       0,
			t1:       1515150445,
			stepTime: 30,
			result:   []byte{0x0, 0x0, 0x0, 0x0, 0x3, 0x2, 0xa5, 0x36},
		},
	}

	for _, c := range tests {
		msg := getMessage(c.t1, c.t0, c.stepTime)
		if !bytes.Equal(msg, c.result) {
			t.Fatalf("expected %#v but got %#v\n", c.result, msg)
		}
	}
}
