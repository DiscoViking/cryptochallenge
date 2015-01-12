package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
)

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
