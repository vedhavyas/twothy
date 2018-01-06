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

func Test_getOffset(t *testing.T) {
	tests := []struct {
		bytes  []byte
		result int
	}{
		{
			bytes:  []byte{0x5a},
			result: 10,
		},

		{
			bytes:  []byte{0x12, 0x12},
			result: 2,
		},

		{
			bytes:  []byte{0xaf},
			result: 15,
		},
	}

	for _, c := range tests {
		r := getOffset(c.bytes)
		if r != c.result {
			t.Fatalf("Expected %d but got %d", c.result, r)
		}
	}
}

func Test_pow(t *testing.T) {
	tests := []struct {
		x int
		y int
		r int
	}{
		{
			x: 0,
			y: 100,
			r: 0,
		},

		{
			x: 100,
			y: 0,
			r: 1,
		},

		{
			x: 2,
			y: 1,
			r: 2,
		},

		{
			x: 10,
			y: 6,
			r: 1000000,
		},
	}

	for _, c := range tests {
		r := pow(c.x, c.y)
		if c.r != r {
			t.Fatalf("expected %d but got %d", c.r, r)
		}
	}
}
