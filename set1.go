package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"
)

// Crypto Challenge Set 1
//
// This is the qualifying set.
// We picked the exercises in it to ramp developers up gradually into coding cryptography,
// but also to verify that we were working with people who were ready to write code.
//
// This set is relatively easy. With one exception, most of these exercises should take only a couple minutes.
// But don't beat yourself up if it takes longer than that. It took Alex two weeks to get through the set!
//
// If you've written any crypto code in the past, you're going to feel like skipping a lot of this.
// Don't skip them. At least two of them (we won't say which) are important stepping stones to later attacks.
var set1 = []func() bool{
	set1_1,
	set1_2,
	set1_3,
	set1_4,
}

// Set 1 Challenge 1
//
// Convert hex to base64
// The string:
//
// 49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d
// Should produce:
//
// SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t
// So go ahead and make that happen. You'll need to use this code for the rest of the exercises.
func set1_1() bool {
	input := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	expected := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	got, err := hexTo64(input)
	if err != nil {
		fmt.Print("Failed to run: ", err)
		return false
	}

	if got != expected {
		fmt.Printf("Wrong answer.\nExpected: %s\nGot: %s\n", expected, got)
		return false
	}
	return true
}

// Set 1 Challenge 2
//
// Fixed XOR
// Write a function that takes two equal-length buffers and produces their XOR combination.
//
// If your function works properly, then when you feed it the string:
//
// 1c0111001f010100061a024b53535009181c
// ... after hex decoding, and when XOR'd against:
//
// 686974207468652062756c6c277320657965
// ... should produce:
//
// 746865206b696420646f6e277420706c6179
func set1_2() bool {
	input_a := "1c0111001f010100061a024b53535009181c"
	input_b := "686974207468652062756c6c277320657965"
	expected := "746865206b696420646f6e277420706c6179"

	hex_a, err := hex.DecodeString(input_a)
	if err != nil {
		fmt.Print("Failed to decode a: ", err)
		return false
	}

	hex_b, err := hex.DecodeString(input_b)
	if err != nil {
		fmt.Print("Failed to decode b: ", err)
		return false
	}

	raw, err := xorBytes(hex_a, hex_b)
	if err != nil {
		fmt.Print("Failed to run xor: ", err)
		return false
	}

	got := hex.EncodeToString(raw)

	if got != expected {
		fmt.Printf("Wrong answer.\nExpected: %s\nGot: %s\n", expected, got)
		return false
	}
	return true
}

// Set 1 Challenge 3
//
// Single-byte XOR cipher
// The hex encoded string:
//
// 1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736
// ... has been XOR'd against a single character. Find the key, decrypt the message.
//
// You can do this by hand. But don't: write code to do it for you.
//
// How? Devise some method for "scoring" a piece of English plaintext.
// Character frequency is a good metric.
// Evaluate each output and choose the one with the best score.
func set1_3() bool {
	ciphertext := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	cipherbytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		fmt.Println("Failed to decode hex: ", err)
		return false
	}

	bestBytes, _, err := crackSingleByteXor(cipherbytes)
	if err != nil {
		fmt.Println("Failed to crack: ", err)
		return false
	}

	plaintext := string(bestBytes)
	fmt.Println("Plaintext: ", plaintext)

	return true
}

// Set 1 Challenge 4
//
// Detect single-character XOR
// One of the 60-character strings in this file has been encrypted by single-character XOR.
//
// Find it.
//
// (Your code from #3 should help.)
func set1_4() bool {
	file, err := os.Open("./data/4.txt")
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return false
	}

	s := bufio.NewScanner(file)

	var bestBytes []byte
	var bestScore float64 = math.MaxFloat64
	lnum := 0
	bestLine := 0
	for s.Scan() {
		line := s.Text()
		cipherbytes, err := hex.DecodeString(line)
		if err != nil {
			fmt.Println("Failed to decode hex: ", err)
			return false
		}

		lnum++
		candidate, score, err := crackSingleByteXor(cipherbytes)
		if err != nil {
			fmt.Println("Failed to crack: ", err)
			return false
		}

		if score < bestScore {
			bestLine = lnum
			bestBytes = candidate
			bestScore = score
		}
	}

	plaintext := string(bestBytes)
	fmt.Println("Ciphertext was on line ", bestLine)
	fmt.Println("Plaintext: ", plaintext)

	return true
}
