// Test file for pw.go
//
// Test data blatantly stolen and modified from
// https://code.google.com/p/go/source/browse/scrypt/scrypt_test.go?repo=crypto
//
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package pw

import (
	"crypto/rand"
	"io"
	"testing"
)

type testPw struct {
	pw     string // password
	hmk    []byte // HMAC Key
	salt   []byte // salt
	h      []byte // hash to check
	output bool   // Output should be true
}

var emptyHash = []byte{}
var emptyString = ""

var passNorm = "password"
var passShort = "pw"
var passLong = "this is a long \000 password"

var hashNorm = []byte{
	0x70, 0x23, 0xbd, 0xcb, 0x3a, 0xfd, 0x73, 0x48,
	0x46, 0x1c, 0x06, 0xcd, 0x81, 0xfd, 0x38, 0xeb,
	0xfd, 0xa8, 0xfb, 0xba, 0x90, 0x4f, 0x8e, 0x3e,
	0xa9, 0xb5, 0x43, 0xf6, 0x54, 0x5d, 0xa1, 0xf2,
}
var hashShort = []byte{
	0x01,
}
var hashLong = []byte{
	0xc3, 0xf1, 0x82, 0xee, 0x2d, 0xec, 0x84, 0x6e,
	0x70, 0xa6, 0x94, 0x2f, 0xb5, 0x29, 0x98, 0x5a,
	0x3a, 0x09, 0x76, 0x5e, 0xf0, 0x4c, 0x61, 0x29,
	0x23, 0xb1, 0x7f, 0x18, 0x55, 0x5a, 0x37, 0x07,
	0x6d, 0xeb, 0x2b, 0x98, 0x30, 0xd6, 0x9d, 0xe5,
	0x49, 0x26, 0x51, 0xe4, 0x50, 0x6a, 0xe5, 0x77,
	0x6d, 0x96, 0xd4, 0x0f, 0x67, 0xaa, 0xee, 0x37,
	0xe1, 0x77, 0x7b, 0x8a, 0xd5, 0xc3, 0x11, 0x14,
	0x32, 0xbb, 0x3b, 0x6f, 0x7e, 0x12, 0x64, 0x40,
	0x18, 0x79, 0xe6, 0x41, 0xae,
}

var hmacNorm = hashNorm
var hmacShort = hashShort
var hmacLong = hashLong

var saltNorm = hashNorm
var saltShort = hashShort
var saltLong = hashLong

var good = []testPw{
	{
		passNorm,
		hmacNorm,
		saltNorm,
		[]byte{
			0x09, 0x8f, 0xd0, 0x1f, 0x55, 0x10, 0xb4, 0x00,
			0x6c, 0x1f, 0xbd, 0xd8, 0xeb, 0x32, 0x99, 0x35,
			0x72, 0x70, 0x1c, 0x53, 0x1f, 0xb9, 0x05, 0x82,
			0xc3, 0x4c, 0xb8, 0x2a, 0xe9, 0xef, 0xe3, 0xf7,
		},
		true,
	},
	{
		passNorm,
		hmacNorm,
		saltShort,
		[]byte{
			0x06, 0x57, 0xcc, 0x93, 0x26, 0xb5, 0xe3, 0xea,
			0x1f, 0x33, 0xb1, 0x62, 0x2b, 0x57, 0xdd, 0xdd,
			0xc1, 0xaf, 0x35, 0x47, 0xc2, 0x45, 0x29, 0xb7,
			0xb5, 0xcd, 0x1b, 0xca, 0x20, 0x3f, 0x6d, 0xcb,
		},
		true,
	},
	{
		passNorm,
		hmacNorm,
		saltLong,
		[]byte{
			0x18, 0x1f, 0x7e, 0xf9, 0x62, 0x3f, 0x50, 0x51,
			0x9e, 0x29, 0x38, 0x7f, 0x77, 0xd6, 0x70, 0x93,
			0xd6, 0xc0, 0x78, 0xc7, 0x7d, 0x5d, 0x25, 0x18,
			0xba, 0xe9, 0x18, 0x7a, 0xc9, 0xa1, 0x26, 0x6c,
		},
		true,
	},
	{
		passShort,
		hmacNorm,
		saltNorm,
		[]byte{
			0xcc, 0x02, 0xd4, 0x6d, 0x29, 0x1f, 0xd6, 0x96,
			0x04, 0xf7, 0x8a, 0xc0, 0x81, 0xe1, 0x4b, 0x57,
			0x47, 0x75, 0x3a, 0x70, 0x54, 0xbe, 0x14, 0x68,
			0x0e, 0xa5, 0xc7, 0x66, 0xc5, 0x75, 0x02, 0xde,
		},
		true,
	},
	{
		passShort,
		hmacNorm,
		saltShort,
		[]byte{
			0x11, 0xd9, 0x99, 0xfc, 0xe8, 0x97, 0xf4, 0x3a,
			0xc9, 0x4a, 0x71, 0x4d, 0xc4, 0xcd, 0x60, 0xcf,
			0x07, 0x86, 0xef, 0xf1, 0x1a, 0xeb, 0xec, 0x58,
			0xe2, 0xd6, 0xf4, 0x61, 0xb8, 0xad, 0x46, 0x33,
		},
		true,
	},
	{
		passShort,
		hmacNorm,
		saltLong,
		[]byte{
			0x96, 0xda, 0xfc, 0x58, 0x4d, 0x03, 0xaf, 0xcf,
			0x3a, 0x2c, 0x9f, 0xf1, 0x24, 0xbf, 0x0b, 0x95,
			0xea, 0x28, 0x0a, 0x8b, 0x17, 0xb2, 0xa6, 0x01,
			0x8e, 0x18, 0x0b, 0x37, 0xb9, 0x24, 0x95, 0x7b,
		},
		true,
	},
	{
		passLong,
		hmacNorm,
		saltNorm,
		[]byte{
			0x62, 0x94, 0x0d, 0xd0, 0xa5, 0xfd, 0xca, 0xf7,
			0xa5, 0x3e, 0xdb, 0x8d, 0x46, 0x5a, 0x07, 0xef,
			0xc2, 0x69, 0xc5, 0x88, 0xd3, 0x78, 0x8b, 0x52,
			0x24, 0x48, 0x7a, 0xc1, 0xac, 0xec, 0x28, 0x5e,
		},
		true,
	},
	{
		passLong,
		hmacNorm,
		saltShort,
		[]byte{
			0x57, 0x35, 0x4a, 0x63, 0x8f, 0xfc, 0x89, 0xba,
			0x1c, 0x80, 0xdd, 0x9d, 0x6b, 0xd5, 0x2e, 0xd5,
			0x8b, 0xb4, 0xe4, 0x00, 0x00, 0x1e, 0xf3, 0xe3,
			0xcb, 0xf7, 0xbe, 0xd9, 0x76, 0x8c, 0x75, 0x43,
		},
		true,
	},
	{
		passLong,
		hmacNorm,
		saltLong,
		[]byte{
			0xdf, 0x9a, 0xd1, 0x14, 0x00, 0xb2, 0x67, 0xaf,
			0xa0, 0xc2, 0x7b, 0xb4, 0x15, 0xce, 0xe6, 0xfc,
			0x28, 0x69, 0x07, 0x7a, 0x55, 0x47, 0xec, 0xf8,
			0xc1, 0x30, 0x46, 0x43, 0x7c, 0x2a, 0x91, 0xbb,
		},
		true,
	},
	{
		passNorm,
		hmacShort,
		saltNorm,
		[]byte{
			0x0b, 0x79, 0x7d, 0xbf, 0xe8, 0x4e, 0x08, 0xdc,
			0x8d, 0x1f, 0x50, 0x8e, 0x9f, 0xa3, 0x6a, 0x59,
			0xd5, 0x24, 0x74, 0x43, 0xa8, 0x73, 0x24, 0x08,
			0x34, 0x2e, 0xd1, 0x12, 0x59, 0xe1, 0xf3, 0x99,
		},
		true,
	},
	{
		passNorm,
		hmacShort,
		saltShort,
		[]byte{
			0xa0, 0x80, 0x6b, 0x62, 0xe7, 0x74, 0xc4, 0x1d,
			0x9e, 0x86, 0xf2, 0x3e, 0x7f, 0x61, 0x10, 0xaf,
			0x06, 0xb0, 0x06, 0xbd, 0x18, 0x3b, 0x08, 0xf8,
			0x2a, 0xfe, 0xa4, 0xea, 0xaa, 0xeb, 0xa9, 0xba,
		},
		true,
	},
	{
		passNorm,
		hmacShort,
		saltLong,
		[]byte{
			0x7b, 0x5f, 0xb5, 0x57, 0x6e, 0xa8, 0xfd, 0x15,
			0x95, 0x41, 0xde, 0xa6, 0x7f, 0x4b, 0x54, 0x39,
			0xd6, 0xa3, 0x5e, 0x5c, 0xc6, 0x29, 0x34, 0x02,
			0x31, 0x84, 0x6b, 0xfd, 0x82, 0xf6, 0x15, 0x1c,
		},
		true,
	},
	{
		passShort,
		hmacShort,
		saltNorm,
		[]byte{
			0xbc, 0x99, 0x0f, 0x51, 0xb0, 0xb6, 0x0f, 0x1f,
			0x8c, 0xf0, 0xd8, 0x90, 0x8d, 0xb9, 0x00, 0x53,
			0x7c, 0xcb, 0x8a, 0xf7, 0x18, 0xd2, 0x0f, 0x6b,
			0xf5, 0xc0, 0xe3, 0x27, 0x17, 0x9b, 0x1c, 0x95,
		},
		true,
	},
	{
		passShort,
		hmacShort,
		saltShort,
		[]byte{
			0xb5, 0x3b, 0x7c, 0x50, 0xf2, 0x7f, 0x4d, 0x31,
			0x28, 0x30, 0x0b, 0xb0, 0x02, 0x57, 0x56, 0x60,
			0x0f, 0x67, 0x83, 0x77, 0x89, 0x33, 0xd2, 0x6a,
			0x92, 0x98, 0xcf, 0xed, 0xf6, 0x53, 0xea, 0x94,
		},
		true,
	},
	{
		passShort,
		hmacShort,
		saltLong,
		[]byte{
			0x91, 0xb4, 0x55, 0x5c, 0x6d, 0x30, 0x0a, 0x30,
			0xe5, 0xf6, 0x97, 0xa0, 0xbc, 0x3a, 0xb9, 0x9a,
			0x67, 0x86, 0x9b, 0xa3, 0x5e, 0x8c, 0x77, 0xc9,
			0xe5, 0xd9, 0xc4, 0x2b, 0x62, 0x02, 0x76, 0x20,
		},
		true,
	},
	{
		passLong,
		hmacShort,
		saltNorm,
		[]byte{
			0x11, 0x2c, 0x99, 0xa6, 0x96, 0x24, 0x4f, 0xec,
			0xf6, 0x40, 0x18, 0xdb, 0xbe, 0x6c, 0x0a, 0xce,
			0xf1, 0x93, 0xf1, 0xc6, 0xc7, 0x13, 0xd4, 0x14,
			0x8a, 0x26, 0xe2, 0xa4, 0x14, 0xfb, 0x7d, 0x52,
		},
		true,
	},
	{
		passLong,
		hmacShort,
		saltShort,
		[]byte{
			0xb0, 0x74, 0xc7, 0x16, 0x6b, 0x1c, 0xdc, 0x8a,
			0xb6, 0x81, 0x4c, 0x58, 0xa2, 0x7f, 0x37, 0x59,
			0xff, 0x1a, 0x13, 0x0b, 0x48, 0x26, 0x0f, 0x0e,
			0x42, 0x71, 0x6b, 0x42, 0x5f, 0x1f, 0x4a, 0x89,
		},
		true,
	},
	{
		passLong,
		hmacShort,
		saltLong,
		[]byte{
			0x6b, 0x51, 0x79, 0xda, 0x82, 0x52, 0x58, 0xcb,
			0x65, 0xb9, 0x96, 0xf6, 0x15, 0x25, 0xb1, 0x37,
			0x08, 0xaa, 0xfe, 0x56, 0xd6, 0x19, 0x2a, 0xaf,
			0x78, 0xdf, 0xea, 0x27, 0xa4, 0x3b, 0x8f, 0x2f,
		},
		true,
	},
	{
		passNorm,
		hmacLong,
		saltNorm,
		[]byte{
			0xa2, 0x87, 0xa8, 0x2d, 0x2a, 0x65, 0xad, 0x44,
			0xd1, 0xaa, 0x64, 0x72, 0x32, 0x8b, 0x9a, 0x51,
			0x2a, 0xcd, 0x36, 0xa2, 0x9c, 0xa8, 0x9f, 0x16,
			0xdd, 0x25, 0xfc, 0xdc, 0xca, 0xa5, 0x44, 0xee,
		},
		true,
	},
	{
		passNorm,
		hmacLong,
		saltShort,
		[]byte{
			0xf7, 0xe8, 0xa6, 0xd8, 0x87, 0x0f, 0xf4, 0xc5,
			0x60, 0xde, 0xca, 0x18, 0x66, 0x66, 0x28, 0x68,
			0x71, 0xd3, 0x76, 0xd2, 0x64, 0xe3, 0x3b, 0x08,
			0x1f, 0x73, 0xbe, 0xfe, 0x95, 0x63, 0xe1, 0x9c,
		},
		true,
	},
	{
		passNorm,
		hmacLong,
		saltLong,
		[]byte{
			0x9b, 0x50, 0x56, 0xf2, 0xcb, 0x9c, 0x8b, 0x99,
			0x2d, 0x47, 0xca, 0xa6, 0xb2, 0x80, 0x1f, 0xe4,
			0xd8, 0xec, 0xa9, 0x7d, 0xeb, 0xcd, 0x98, 0x9a,
			0x33, 0xb1, 0xfa, 0x7c, 0x6d, 0xfe, 0x91, 0xbe,
		},
		true,
	},
	{
		passShort,
		hmacLong,
		saltNorm,
		[]byte{
			0x8e, 0xfb, 0x47, 0x3f, 0xde, 0x18, 0xdf, 0xbb,
			0x37, 0x35, 0x12, 0xc4, 0xd3, 0x8a, 0x80, 0x2a,
			0x05, 0x45, 0x06, 0x9e, 0x8d, 0x87, 0x35, 0x87,
			0x92, 0xe1, 0x8b, 0xec, 0xff, 0xc7, 0xf5, 0x2b,
		},
		true,
	},
	{
		passShort,
		hmacLong,
		saltShort,
		[]byte{
			0x41, 0x41, 0xe7, 0xa2, 0xf1, 0x59, 0x42, 0xc4,
			0xab, 0x8f, 0x0e, 0xbe, 0x63, 0xb2, 0x70, 0x2d,
			0x56, 0xf5, 0x91, 0x96, 0x11, 0x8c, 0xfc, 0xee,
			0xc5, 0xf4, 0x0f, 0x87, 0xab, 0xf9, 0x77, 0xdf,
		},
		true,
	},
	{
		passShort,
		hmacLong,
		saltLong,
		[]byte{
			0x9d, 0xd3, 0x50, 0x5a, 0x33, 0xb6, 0x50, 0xc2,
			0x31, 0x33, 0x26, 0x78, 0x0f, 0xf8, 0xed, 0x06,
			0x1c, 0x75, 0x91, 0x79, 0x42, 0x85, 0x6c, 0xcd,
			0x94, 0x27, 0x6f, 0xfc, 0x30, 0x50, 0x93, 0xe2,
		},
		true,
	},
	{
		passLong,
		hmacLong,
		saltNorm,
		[]byte{
			0x67, 0xb9, 0xf2, 0x6e, 0x95, 0x90, 0x7e, 0x5e,
			0x90, 0x05, 0x21, 0x4d, 0xf3, 0x28, 0xb8, 0x02,
			0x83, 0x41, 0x50, 0x69, 0x82, 0x30, 0x74, 0x68,
			0xb1, 0x92, 0x93, 0x04, 0x65, 0x0c, 0x64, 0xad,
		},
		true,
	},
	{
		passLong,
		hmacLong,
		saltShort,
		[]byte{
			0xd1, 0x6a, 0x14, 0xfb, 0x88, 0x36, 0xf7, 0x35,
			0xad, 0xc1, 0x07, 0xae, 0x7b, 0x9d, 0x87, 0xef,
			0x1c, 0x96, 0xa7, 0x2d, 0x2c, 0x6f, 0x71, 0xfe,
			0x7f, 0x98, 0x0a, 0xba, 0xae, 0xc3, 0x6e, 0xdb,
		},
		true,
	},
	{
		passLong,
		hmacLong,
		saltLong,
		[]byte{
			0xbf, 0xda, 0x0e, 0x92, 0xac, 0x06, 0x08, 0xc6,
			0xb9, 0x7e, 0x16, 0x95, 0x7e, 0x54, 0xda, 0x8a,
			0x85, 0x0d, 0xe4, 0xf7, 0x34, 0x82, 0x84, 0x2f,
			0xfe, 0x9c, 0x74, 0x4c, 0x43, 0xbe, 0xcc, 0x5f,
		},
		true,
	},
	{
		emptyString,
		emptyHash,
		emptyHash,
		[]byte{
			0xd8, 0x5f, 0x86, 0x7d, 0x5e, 0x12, 0xd9, 0xf5,
			0xcf, 0xca, 0x55, 0x9c, 0x65, 0x85, 0x8e, 0x3c,
			0x22, 0x29, 0x60, 0x81, 0x5f, 0x0a, 0xfa, 0x44,
			0x1c, 0xdb, 0xc4, 0x9c, 0xab, 0xa5, 0xd3, 0xa7,
		},
		true,
	},
}

var bad = []testPw{
	{emptyString, emptyHash, emptyHash, emptyHash, true}, // No input, No output
	{emptyString, emptyHash, emptyHash, hashNorm, true},  // No input, Wrong output
	{passNorm, hmacNorm, saltNorm, emptyHash, true},      // Normal input, No output
	{passNorm, hmacNorm, saltNorm, hashNorm, true},       // Normal input, Wrong output
}

func TestSet(t *testing.T) {
	if testing.Short() {
		good = good[13:14]
	}
	id := New()
	for _, v := range good {
		id.Pass, id.Hmac, id.Salt = v.pw, v.hmk, v.salt
		if err := id.Set(); err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
	}
	// Invalid Scrypt variables
	id.N, id.R, id.P = 0, 0, 0
	if err := id.Set(); err == nil {
		t.Errorf("expected err, got nil")
	}
}

func TestCheck(t *testing.T) {
	id := New()
	for i, v := range good {
		id.Pass, id.Hmac, id.Salt, id.Hash = v.pw, v.hmk, v.salt, v.h
		chk, err := id.Verify()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
		if chk != true {
			t.Errorf("%d: expected %t, got %t", i, v.output, chk)
		}
	}
	for i, v := range bad {
		id.Pass, id.Hmac, id.Salt, id.Hash = v.pw, v.hmk, v.salt, v.h
		chk, err := id.Verify()
		if err == nil {
			t.Errorf("%d: expected error, got nil, function returned %t", i, chk)
		}
	}
	// Invalid Scrypt variables
	id.N, id.R, id.P = 0, 0, 0
	if _, err := id.Verify(); err == nil {
		t.Errorf("expected err, got nil")
	}
}

func TestSetAndCheck(t *testing.T) {
	id := New()
	for i, v := range good {
		id.Pass, id.Hmac, id.Salt = v.pw, v.hmk, v.salt
		if err := id.Set(); err != nil {
			t.Errorf("got unexpected error: %v", err)
		}
		chk, err := id.Verify()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
		if chk != true {
			t.Errorf("%d: expected %t, got %t", i, v.output, chk)
		}
	}
}

func TestRandomSetAndCheck(t *testing.T) {
	id := New()
	if err := id.randSalt(); err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	tmpPass := id.Salt
	id.Pass = string(tmpPass)
	if err := id.Set(); err != nil {
		t.Errorf("got unexpected error: %v", err)
	}
	chk, err := id.Verify()
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	if chk != true {
		t.Errorf("expected %t, got %t", true, chk)
	}
}

func TestRandSalt(t *testing.T) {
	id := New()
	randSrc = io.LimitReader(rand.Reader, 0)
	if err := id.Set(); err != nil && err != io.EOF {
		t.Errorf("got unexpected error: %v", err)
	}
	if _, err := id.Verify(); err == nil {
		t.Errorf("expected err, got nil")
	}
}

func TestScrypt(t *testing.T) {
	id := New()
	// Invalid Scrypt variables
	id.N, id.R, id.P = 0, 0, 0
	if err := id.doHash(); err == nil {
		t.Errorf("expected err, got nil")
	}
}

func BenchmarkSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Pass, id.Hmac, id.Salt = good[1].pw, good[1].hmk, good[1].salt
		if err := id.Set(); err == nil {
			b.Errorf("expected err, got nil")
		}
	}
}

func BenchmarkCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Pass, id.Hmac, id.Salt, id.Hash = good[1].pw, good[1].hmk, good[1].salt, good[1].h
		if _, err := id.Verify(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}

func BenchmarkSetAndCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Pass, id.Hmac, id.Salt = good[1].pw, good[1].hmk, good[1].salt
		if err := id.Set(); err == nil {
			b.Errorf("expected %v, got nil", err)
		}
		if _, err := id.Verify(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}
