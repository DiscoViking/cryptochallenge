package main

// General utility functions.

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
)

func hexTo64(s string) (string, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func xorBytes(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, errors.New("lengths of byte arrays do not match")
	}

	out := make([]byte, len(a))

	for i := 0; i < len(a); i++ {
		out[i] = a[i] ^ b[i]
	}

	return out, nil
}
