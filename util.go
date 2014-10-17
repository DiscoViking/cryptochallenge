package main

// General utility functions.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"math"
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

func crackSingleByteXor(cipherbytes []byte) ([]byte, float64, error) {
	var bestBytes []byte
	var bestScore float64 = math.MaxFloat64
	key := make([]byte, len(cipherbytes))

	for b := byte(0); b < 255; b++ {
		for i := range key {
			key[i] = b
		}

		candidate, err := xorBytes(cipherbytes, key)
		if err != nil {
			return nil, 0, err
		}

		score := englishness(bytes.ToLower(candidate))
		if score < bestScore {
			bestBytes = candidate
			bestScore = score
		}
	}

	return bestBytes, bestScore, nil
}
