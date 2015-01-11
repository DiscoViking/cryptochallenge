package main

// General utility functions.

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
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

// Repeating key Xor
func repeatingKeyXor(key, input []byte) ([]byte, error) {
	if len(key) == 0 {
		return nil, errors.New("zero-length key")
	}
	out := make([]byte, len(input))
	keyIx := 0

	for i, b := range input {
		out[i] = key[keyIx] ^ b
		if keyIx++; keyIx >= len(key) {
			keyIx = 0
		}
	}

	return out, nil
}

// Cracks a single byte Xor cipher.
func crackSingleByteXor(cipherbytes []byte) ([]byte, byte, float64, error) {
	var bestBytes []byte
	var bestScore float64 = 0
	var bestKey byte = 0

	for b := byte(0); b < 255; b++ {
		candidate, err := repeatingKeyXor([]byte{b}, cipherbytes)
		if err != nil {
			return nil, 0, 0, err
		}

		score := englishness(bytes.ToLower(candidate))
		//fmt.Println(hex.EncodeToString([]byte{b}), ":", score, ": ", string(candidate))
		if score > bestScore {
			bestBytes = candidate
			bestScore = score
			bestKey = b
		}
	}

	fmt.Println("Best :", bestScore, ": ", string(bestBytes))

	return bestBytes, bestKey, bestScore, nil
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

// Finds the most likely key of a given length for a repeating key Xor.
// Returns plainbytes, key, score.
func findLikelyXorKey(cipherbytes []byte, n int) ([]byte, []byte, error) {
	// Now transpose by the assumed key size.
	cols := transposeBytes(cipherbytes, uint64(n))
	plainCols := [][]byte{}
	key := make([]byte, 0, n)

	// Crack each column individually.
	for _, col := range cols {
		fmt.Printf("\n\nCracking: %v\n\n", hex.EncodeToString(col))
		plain, b, _, err := crackSingleByteXor(col)
		if err != nil {
			return nil, nil, err
		}

		plainCols = append(plainCols, plain)
		key = append(key, b)
	}

	// And recombine the solved columns.
	plainbytes := combineColumns(plainCols)

	return plainbytes, key, nil
}

// Cracks a repeating key Xor cipher.
// Returns the plaintext bytes, the key and maybe an error.
func crackRepeatingKeyXor(cipherbytes []byte) ([]byte, []byte, error) {
	minKeySize := 2
	maxKeySize := 40
	bestKeySize := 0
	var bestScore float64 = 8 // Worst possible score.

	// Firstly, calculate the keysize by comparing blocks of bytes.
	for keySize := minKeySize; keySize <= maxKeySize; keySize++ {
		if len(cipherbytes) < 4*keySize {
			break
		}

		// Calculate the hamming distance between the first two blocks of keysize bytes.
		a := cipherbytes[0:keySize]
		b := cipherbytes[keySize : 2*keySize]
		c := cipherbytes[2*keySize : 3*keySize]
		d := cipherbytes[3*keySize : 4*keySize]
		var total uint64 = 0

		h, err := hammingDistance(a, b)
		if err != nil {
			return nil, nil, err
		}
		total += h

		h, err = hammingDistance(b, c)
		if err != nil {
			return nil, nil, err
		}
		total += h

		h, err = hammingDistance(c, d)
		if err != nil {
			return nil, nil, err
		}
		total += h

		total /= 3

		normalized := float64(total) / float64(keySize)
		fmt.Println(keySize, " - ", normalized)

		// Keep track of the best score.
		if normalized < bestScore {
			bestScore = normalized
			bestKeySize = keySize
		}
	}
	bestKeySize = 29

	plainbytes, key, err := findLikelyXorKey(cipherbytes, bestKeySize)
	if err != nil {
		return nil, nil, err
	}
	return plainbytes, key, nil
}
