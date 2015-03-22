// Test file for pw.go
// Blatantly stolen and modified from https://code.google.com/p/go/source/browse/scrypt/scrypt_test.go?repo=crypto
//
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pw

import "testing"

type testPw struct {
	hmk    []byte // HMAC Key
	pw     string // password
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
		hmacNorm,
		passNorm,
		saltNorm,
		[]byte{
			0x60, 0x72, 0x94, 0xca, 0xef, 0x28, 0x4b, 0xb3,
			0x15, 0x8c, 0x2c, 0x9, 0xf9, 0x4e, 0xf2, 0x43,
			0xda, 0x81, 0x14, 0x6c, 0xfd, 0xcf, 0xef, 0xfd,
			0xe9, 0x75, 0xc0, 0x72, 0x30, 0xe2, 0x73, 0x22,
		},
		true,
	},
	{
		hmacNorm,
		passNorm,
		saltShort,
		[]byte{
			0x1d, 0x30, 0x6c, 0x45, 0x27, 0xed, 0x8d, 0x8a,
			0x8c, 0x54, 0x87, 0x44, 0x1e, 0x61, 0x63, 0x60,
			0x94, 0xd, 0x94, 0xb6, 0x4f, 0x61, 0xf9, 0x36,
			0x5a, 0x5d, 0x3f, 0x8a, 0xba, 0xc7, 0x69, 0xd5,
		},
		true,
	},
	{
		hmacNorm,
		passNorm,
		saltLong,
		[]byte{
			0xcb, 0xb4, 0x8b, 0xb6, 0x65, 0x63, 0xc0, 0x35,
			0x3e, 0x3a, 0x54, 0xa7, 0xdc, 0xed, 0xbd, 0x13,
			0xf, 0x6e, 0x87, 0xd6, 0xd4, 0x8e, 0xb9, 0x7,
			0x8, 0xfc, 0x2b, 0x42, 0xdb, 0x52, 0x83, 0xd3,
		},
		true,
	},
	{
		hmacNorm,
		passShort,
		saltNorm,
		[]byte{
			0x81, 0x79, 0xd6, 0x34, 0x7, 0x2d, 0xd0, 0xf9,
			0xd8, 0xbf, 0xec, 0x70, 0xcc, 0xf8, 0x8a, 0x77,
			0xe6, 0xb, 0xb1, 0x80, 0x8b, 0xfe, 0x9d, 0xf,
			0x19, 0x66, 0x8f, 0x50, 0x72, 0xfe, 0x64, 0x9f,
		},
		true,
	},
	{
		hmacNorm,
		passShort,
		saltShort,
		[]byte{
			0xd, 0xf3, 0xad, 0xa4, 0x32, 0x9c, 0x53, 0xd,
			0x6f, 0xb0, 0x34, 0x3, 0xa2, 0x88, 0x85, 0x48,
			0x12, 0x4, 0xe7, 0x16, 0xb8, 0xff, 0xa1, 0x6,
			0x56, 0x58, 0xd1, 0x61, 0x41, 0x1d, 0x7b, 0x64,
		},
		true,
	},
	{
		hmacNorm,
		passShort,
		saltLong,
		[]byte{
			0xcf, 0xec, 0xe9, 0xa1, 0x7, 0x6, 0x5d, 0x69,
			0xa1, 0x59, 0x59, 0xc2, 0xcd, 0x20, 0xd3, 0xfd,
			0x18, 0xc3, 0x73, 0xbd, 0xf2, 0x62, 0xe0, 0x9c,
			0x9c, 0xcf, 0x41, 0x44, 0x1a, 0x44, 0x84, 0xa4,
		},
		true,
	},
	{
		hmacNorm,
		passLong,
		saltNorm,
		[]byte{
			0x5c, 0xac, 0xd5, 0x1, 0x27, 0xf6, 0x47, 0xe2,
			0x11, 0x32, 0x4b, 0x73, 0x7, 0xcb, 0xff, 0x3b,
			0xe8, 0x2, 0xc4, 0xc6, 0x3f, 0x3, 0x6e, 0x3b,
			0xab, 0xbf, 0x69, 0xee, 0x1, 0xd0, 0x76, 0x98,
		},
		true,
	},
	{
		hmacNorm,
		passLong,
		saltShort,
		[]byte{
			0x7e, 0x78, 0x31, 0xc5, 0x60, 0x1d, 0x21, 0xfe,
			0x23, 0x13, 0x4f, 0x10, 0xa9, 0x6a, 0x8e, 0xf8,
			0x83, 0x14, 0x90, 0xd6, 0x36, 0x42, 0xf, 0xad,
			0x90, 0x65, 0x3f, 0x8e, 0x8c, 0x9e, 0x9c, 0x9d,
		},
		true,
	},
	{
		hmacNorm,
		passLong,
		saltLong,
		[]byte{
			0x9e, 0xed, 0x30, 0x54, 0xc0, 0x49, 0x92, 0xd4,
			0x60, 0xe6, 0xc4, 0xa3, 0xab, 0xff, 0x27, 0x30,
			0xc2, 0x2b, 0x89, 0x42, 0x4, 0x99, 0x7b, 0xfb,
			0xfa, 0x2, 0xa3, 0xc4, 0xff, 0x32, 0x42, 0x53,
		},
		true,
	},
	{
		hmacShort,
		passNorm,
		saltNorm,
		[]byte{
			0x5b, 0xd9, 0xb0, 0x74, 0xc5, 0x69, 0xec, 0xf4,
			0xd3, 0x3f, 0x66, 0xf3, 0x80, 0x1d, 0xa8, 0xb2,
			0xd8, 0x9, 0xfd, 0x3, 0xa9, 0x6e, 0x2c, 0x7d,
			0xbd, 0x24, 0x86, 0x80, 0xec, 0x1c, 0xac, 0x1a,
		},
		true,
	},
	{
		hmacShort,
		passNorm,
		saltShort,
		[]byte{
			0x70, 0x2c, 0x1d, 0x63, 0x55, 0x6d, 0xee, 0xe9,
			0x57, 0xc4, 0x72, 0x21, 0x41, 0x9, 0x9b, 0x4d,
			0xfd, 0xe7, 0x3b, 0x36, 0x1e, 0x53, 0x44, 0xdc,
			0x12, 0xc2, 0x78, 0xd5, 0x1a, 0xe1, 0x36, 0xdf,
		},
		true,
	},
	{
		hmacShort,
		passNorm,
		saltLong,
		[]byte{
			0x70, 0x3a, 0x31, 0x90, 0xb1, 0x22, 0xf2, 0x9,
			0xda, 0x94, 0xe, 0xbc, 0xa6, 0x78, 0xb, 0x70,
			0xd, 0xe6, 0x22, 0x31, 0x6e, 0xa7, 0xc5, 0x10,
			0xc3, 0xbc, 0x28, 0xb8, 0xac, 0x36, 0x89, 0xfb,
		},
		true,
	},
	{
		hmacShort,
		passShort,
		saltNorm,
		[]byte{
			0x6a, 0xb4, 0x75, 0xca, 0x78, 0x17, 0xc2, 0x6c,
			0xd2, 0x34, 0xb0, 0xc3, 0x8, 0x7, 0x32, 0x70,
			0x4a, 0x33, 0x17, 0xfd, 0x4b, 0x39, 0x6d, 0x8d,
			0xa2, 0xa6, 0xa0, 0x1a, 0x2c, 0xd1, 0x69, 0x78,
		},
		true,
	},
	{
		hmacShort,
		passShort,
		saltShort,
		[]byte{
			0x1c, 0x74, 0x77, 0x7, 0x1a, 0x26, 0x3d, 0x96,
			0x8b, 0x16, 0x1a, 0x3f, 0xc3, 0x89, 0xa2, 0x3c,
			0x65, 0x2c, 0xd2, 0x40, 0x7b, 0xbe, 0xdb, 0xa0,
			0xe9, 0x71, 0xce, 0xb9, 0x83, 0xfe, 0x79, 0xb4,
		},
		true,
	},
	{
		hmacShort,
		passShort,
		saltLong,
		[]byte{
			0x4c, 0xd5, 0x2, 0xdd, 0xaa, 0x8c, 0xd8, 0xa0,
			0x40, 0xe7, 0x15, 0xaf, 0x46, 0xd5, 0xf0, 0x8,
			0xf4, 0x6c, 0x59, 0x7d, 0x22, 0x13, 0xf1, 0xbd,
			0xc, 0xb5, 0x67, 0x8, 0x3a, 0x2d, 0x5d, 0x74,
		},
		true,
	},
	{
		hmacShort,
		passLong,
		saltNorm,
		[]byte{
			0x2d, 0x96, 0xec, 0x8e, 0x20, 0x88, 0xc6, 0x6e,
			0x43, 0xd7, 0xab, 0xe5, 0xa5, 0x69, 0xbd, 0x78,
			0xe, 0x58, 0xb0, 0x60, 0x29, 0xf6, 0xbe, 0x80,
			0xde, 0xcb, 0xd2, 0x28, 0x68, 0xec, 0xa7, 0x77,
		},
		true,
	},
	{
		hmacShort,
		passLong,
		saltShort,
		[]byte{
			0x97, 0x66, 0xae, 0x42, 0x3b, 0x8f, 0xaf, 0xf0,
			0x61, 0xd2, 0xbe, 0x93, 0x8, 0xdc, 0x74, 0x5e,
			0xf1, 0x74, 0x37, 0x6f, 0x15, 0x7, 0xfc, 0x7c,
			0xbe, 0x74, 0x38, 0x75, 0xdc, 0xfd, 0xad, 0x5d,
		},
		true,
	},
	{
		hmacShort,
		passLong,
		saltLong,
		[]byte{
			0x39, 0xe0, 0x47, 0xa1, 0xc5, 0xd4, 0x1e, 0x79,
			0xfe, 0x64, 0x4b, 0x8e, 0xf9, 0xee, 0x0, 0x90,
			0x25, 0x4a, 0xed, 0xdf, 0xc, 0x79, 0xe8, 0xbd,
			0x86, 0xe4, 0x3e, 0x8b, 0xcf, 0xa0, 0x10, 0x5,
		},
		true,
	},
	{
		hmacLong,
		passNorm,
		saltNorm,
		[]byte{
			0xb5, 0xa1, 0x87, 0x6d, 0xd, 0xc7, 0xdd, 0xe1,
			0x98, 0x5b, 0x8c, 0x57, 0x2a, 0x17, 0x73, 0x67,
			0x4e, 0x7c, 0x17, 0x5d, 0x48, 0xe6, 0xe9, 0x29,
			0xae, 0x64, 0x21, 0xbb, 0x64, 0x87, 0xfe, 0xd2,
		},
		true,
	},
	{
		hmacLong,
		passNorm,
		saltShort,
		[]byte{
			0x99, 0xd0, 0x2, 0xec, 0x23, 0xb9, 0xb1, 0x8f,
			0xd9, 0x4e, 0xef, 0x0, 0xdd, 0x33, 0x7, 0x1b,
			0xdc, 0x21, 0xdd, 0xde, 0xcb, 0x7a, 0x80, 0x53,
			0x10, 0x47, 0x73, 0x42, 0xeb, 0x22, 0xde, 0x89,
		},
		true,
	},
	{
		hmacLong,
		passNorm,
		saltLong,
		[]byte{
			0xff, 0xd9, 0xd1, 0x5c, 0x7f, 0xae, 0x2c, 0xc2,
			0xd3, 0xd4, 0xe3, 0x4b, 0xb, 0xc1, 0x6b, 0x3f,
			0x3f, 0x97, 0x8f, 0x1, 0xef, 0xca, 0x5f, 0xfa,
			0x7f, 0x37, 0x2a, 0x54, 0x38, 0xc7, 0x91, 0xcc,
		},
		true,
	},
	{
		hmacLong,
		passShort,
		saltNorm,
		[]byte{
			0x75, 0x2d, 0xbc, 0x5a, 0xc, 0x8, 0x58, 0x7f,
			0x13, 0x3b, 0xf6, 0xb7, 0xc1, 0x5, 0xc, 0xf6,
			0xbb, 0x69, 0xca, 0xc8, 0x32, 0x9f, 0x5d, 0xcf,
			0xd7, 0x76, 0x9b, 0x88, 0xa0, 0x97, 0xc4, 0x11,
		},
		true,
	},
	{
		hmacLong,
		passShort,
		saltShort,
		[]byte{
			0x7e, 0xa5, 0x17, 0x3b, 0xa9, 0x3d, 0x90, 0x3f,
			0xd2, 0xe6, 0x4b, 0x84, 0x54, 0xe0, 0xb7, 0x9d,
			0x81, 0x26, 0xa4, 0xc0, 0x42, 0xa, 0x55, 0x2c,
			0x33, 0x7b, 0xd5, 0xe3, 0x38, 0xf4, 0x85, 0xd5,
		},
		true,
	},
	{
		hmacLong,
		passShort,
		saltLong,
		[]byte{
			0x43, 0xcf, 0x66, 0xe1, 0x78, 0xbc, 0x73, 0x61,
			0x1d, 0xda, 0x7a, 0x22, 0xa7, 0xeb, 0xfe, 0xa1,
			0xb8, 0xb7, 0xe5, 0x8a, 0xb0, 0xc8, 0xa4, 0x2f,
			0x40, 0xa8, 0x71, 0x55, 0xf, 0x6a, 0x37, 0x20,
		},
		true,
	},
	{
		hmacLong,
		passLong,
		saltNorm,
		[]byte{
			0xb2, 0xe0, 0xf9, 0x80, 0xf1, 0x70, 0x60, 0xf6,
			0xae, 0xdf, 0x58, 0xf0, 0x5a, 0x38, 0x72, 0x1a,
			0x1d, 0xa9, 0x40, 0x22, 0x5d, 0xba, 0xdf, 0x36,
			0x2b, 0xda, 0x20, 0x1d, 0x34, 0xfa, 0x2d, 0x8c,
		},
		true,
	},
	{
		hmacLong,
		passLong,
		saltShort,
		[]byte{
			0x30, 0x4a, 0x14, 0xdd, 0x43, 0xce, 0x5d, 0x6b,
			0xca, 0xfd, 0xc0, 0x4a, 0x43, 0x1c, 0x58, 0x20,
			0xc6, 0xc9, 0x83, 0x89, 0x8c, 0x53, 0xc1, 0xba,
			0xe4, 0x7e, 0x58, 0x8b, 0x3d, 0x27, 0x17, 0x14,
		},
		true,
	},
	{
		hmacLong,
		passLong,
		saltLong,
		[]byte{
			0x5b, 0xee, 0xa4, 0x58, 0xd6, 0xc, 0x54, 0xbd,
			0x41, 0x2d, 0x25, 0x74, 0x3b, 0xfe, 0x2b, 0xe2,
			0x5, 0xcf, 0xf6, 0xd9, 0xc8, 0x23, 0x46, 0xe9,
			0x17, 0xb9, 0x73, 0x9f, 0xdd, 0x8f, 0x1e, 0x67,
		},
		true,
	},
	{
		emptyHash,
		emptyString,
		emptyHash,
		[]byte{
			0x8e, 0x1f, 0x9d, 0x8a, 0x51, 0x83, 0x3c, 0x69,
			0x86, 0xae, 0xa0, 0xa, 0xe1, 0xf8, 0xe7, 0x82,
			0x17, 0x82, 0xac, 0x7d, 0x94, 0x95, 0x5, 0xf4,
			0x0, 0xbb, 0xb3, 0xe6, 0xcd, 0x44, 0xad, 0x89,
		},
		true,
	},
}

var bad = []testPw{
	{emptyHash, emptyString, emptyHash, emptyHash, true}, // No input, No output
	{emptyHash, emptyString, emptyHash, hashNorm, true},  // No input, Wrong output
	{hmacNorm, passNorm, saltNorm, emptyHash, true},      // Normal input, No output
	{hmacNorm, passNorm, saltNorm, hashNorm, true},       // Normal input, Wrong output
}

func TestCreate(t *testing.T) {
	if testing.Short() {
		good = good[13:14]
	}
	id := New()
	for i, v := range good {
		id.Hmac, id.Pass, id.Salt = v.hmk, v.pw, v.salt
		err := id.Create()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}

func TestCheck(t *testing.T) {
	id := New()
	for i, v := range good {
		id.Hmac, id.Pass, id.Salt, id.Hash = v.hmk, v.pw, v.salt, v.h
		chk, err := id.Check()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
		if chk != true {
			t.Errorf("%d: expected %t, got %t", i, v.output, chk)
		}
	}
	for i, v := range bad {
		id.Hmac, id.Pass, id.Salt, id.Hash = v.hmk, v.pw, v.salt, v.h
		chk, err := id.Check()
		if err == nil {
			t.Errorf("%d: expected error, got nil, function returned %t", i, chk)
		}
	}
}

func TestCreateAndCheck(t *testing.T) {
	id := New()
	for i, v := range good {
		id.Hmac, id.Pass, id.Salt = v.hmk, v.pw, v.salt
		err := id.Create()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
		chk, err := id.Check()
		if err != nil {
			t.Errorf("%d: got unexpected error: %s", i, err)
		}
		if chk != true {
			t.Errorf("%d: expected %t, got %t", i, v.output, chk)
		}
	}
}

func TestRandomCreateAndCheck(t *testing.T) {
	id := New()
	_ = id.randSalt()
	tmpPass := id.Salt
	id.Pass = string(tmpPass)
	err := id.Create()
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	chk, err := id.Check()
	if err != nil {
		t.Errorf("got unexpected error: %s", err)
	}
	if chk != true {
		t.Errorf("expected %t, got %t", true, chk)
	}
}

func BenchmarkCreate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Hmac, id.Pass, id.Salt = good[1].hmk, good[1].pw, good[1].salt
		if err := id.Create(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}

func BenchmarkCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Hmac, id.Pass, id.Salt, id.Hash = good[1].hmk, good[1].pw, good[1].salt, good[1].h
		if _, err := id.Check(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}

func BenchmarkCreateAndCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := New()
		id.Hmac, id.Pass, id.Salt = good[1].hmk, good[1].pw, good[1].salt
		if err := id.Create(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
		if _, err := id.Check(); err != nil {
			b.Errorf("%d: got unexpected error: %s", i, err)
		}
	}
}
