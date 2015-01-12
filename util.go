package main

// General utility functions.

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

// Converts a hex string to a base64 encoded one.
func hexTo64(s string) (string, error) {
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

// Xors two equal-length byte arrays together.
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

// Calculate the Hamming distance between two strings.
// The Hamming distance is just the number of different bits.
// Notice that this is just the number of 1s in A xor B.
func hammingDistance(a, b []byte) (uint64, error) {
	var d uint64 = 0
	x, err := xorBytes(a, b)
	if err != nil {
		return 0, err
	}

	// Now count the 1s in x.
	for _, b := range x {
		for ; b != 0; b = b >> 1 {
			if b&1 == 1 {
				d++
			}
		}
	}

	return d, nil
}

func transposeBytes(b []byte, l uint64) [][]byte {
	// Initialise the colunm array.
	cols := [][]byte{}
	colLen := (uint64(len(b)) / l) + 1

	for i := uint64(0); i < l; i++ {
		col := make([]byte, 0, colLen)
		cols = append(cols, col)
	}

	for offset := uint64(0); offset < l; offset++ {
		for ix := offset; ix < uint64(len(b)); ix += l {
			cols[offset] = append(cols[offset], b[ix])
		}
	}

	return cols
}

func combineColumns(cols [][]byte) []byte {
	if len(cols) == 0 {
		return []byte{}
	}

	colLen := len(cols[0])
	out := make([]byte, 0, colLen*len(cols))

	for i := 0; i < colLen; i++ {
		for _, col := range cols {
			if len(col) <= i {
				continue
			}

			out = append(out, col[i])
		}
	}

	return out
}

// Read a base64 encoded file and decode into bytes.
func readBase64File(f string) ([]byte, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open file: ", err))
	}

	info, err := file.Stat()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to stat file: ", err))
	}

	size := info.Size()
	buf := make([]byte, size)

	n, err := file.Read(buf)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to read file: ", err))
	}

	if int64(n) != size {
		return nil, errors.New(fmt.Sprintf("only read %v bytes out of %v", n, size))
	}

	bytes, err := base64.StdEncoding.DecodeString(string(buf))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to decode base64: %v", err))
	}

	return bytes, nil
}
