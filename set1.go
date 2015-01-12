package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
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
	set1_5,
	set1_6,
	set1_7,
	set1_8,
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

	bestBytes, _, _, err := crackSingleByteXor(cipherbytes)
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
	var bestScore float64 = 0
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
		candidate, _, score, err := crackSingleByteXor(cipherbytes)
		if err != nil {
			fmt.Println("Failed to crack: ", err)
			return false
		}

		if score > bestScore {
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

// Set 1 Challenge 5
//
// Implement repeating-key XOR
// Here is the opening stanza of an important work of the English language:
//
// Burning 'em, if you ain't quick and nimble
// I go crazy when I hear a cymbal
// Encrypt it, under the key "ICE", using repeating-key XOR.
//
// In repeating-key XOR, you'll sequentially apply each byte of the key;
// the first byte of plaintext will be XOR'd against I, the next C, the next E,
// then I again for the 4th byte, and so on.
//
// It should come out to:
//
// 0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272
// a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f
// Encrypt a bunch of stuff using your repeating-key XOR function.
// Encrypt your mail. Encrypt your password file. Your .sig file.
// Get a feel for it. I promise, we aren't wasting your time with this.
func set1_5() bool {
	plaintext := `Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`
	plainbytes := []byte(plaintext)
	key := []byte("ICE")
	expected := `0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f`

	cipherbytes, err := repeatingKeyXor(key, plainbytes)
	if err != nil {
		fmt.Println("Failed to encrypt: ", err)
		return false
	}

	ciphertext := hex.EncodeToString(cipherbytes)

	if ciphertext != expected {
		fmt.Println("Wrong output: ", string(ciphertext))
		return false
	}

	return true
}

// Set 1 Challenge 6
//
// Break repeating-key XOR
// It is officially on, now.
// This challenge isn't conceptually hard, but it involves actual error-prone coding. The other challenges in this set are there to bring you up to speed. This one is there to qualify you. If you can do this one, you're probably just fine up to Set 6.
//
// There's a file here. It's been base64'd after being encrypted with repeating-key XOR.
//
// Decrypt it.
func set1_6() bool {
	cipherbytes, err := readBase64File("./data/6.txt")
	if err != nil {
		fmt.Print(err.Error())
		return false
	}

	plainbytes, key, err := crackRepeatingKeyXor(cipherbytes)
	if err != nil {
		fmt.Println("Failed to crack: ", err)
		return false
	}

	fmt.Println("Found key: ", string(key))
	fmt.Println("Plaintext: ", string(plainbytes))

	return true
}

// Set 1 Challenge 7
// AES in ECB mode
// The Base64-encoded content in this file has been encrypted via AES-128 in ECB mode under the key
//
// "YELLOW SUBMARINE".
// (case-sensitive, without the quotes; exactly 16 characters; I like "YELLOW SUBMARINE" because it's exactly 16 bytes long, and now you do too).
//
// Decrypt it. You know the key, after all.
//
// Easiest way: use OpenSSL::Cipher and give it AES-128-ECB as the cipher.
func set1_7() bool {
	key := []byte("YELLOW SUBMARINE")
	cipherbytes, err := readBase64File("./data/7.txt")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	plainbytes, err := decryptAesEcb(cipherbytes, key)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	plaintext := string(plainbytes)
	fmt.Println("Plaintext:", plaintext)

	return true
}

// Set 1 Challenge 8
// Detect AES in ECB mode
// In this file are a bunch of hex-encoded ciphertexts.
//
// One of them has been encrypted with ECB.
//
// Detect it.
//
// Remember that the problem with ECB is that it is stateless and deterministic;
// the same 16 byte plaintext block will always produce the same 16 byte ciphertext.
func set1_8() bool {
	file, err := os.Open("./data/8.txt")
	if err != nil {
		fmt.Println("Failed to open file: ", err)
		return false
	}

	s := bufio.NewScanner(file)

	lnum := 0
	bestLine := 0
	var bestScore float64 = 16 * 8 // Worst possible score.
	for s.Scan() {
		lnum++
		line := s.Text()
		bytes, err := hex.DecodeString(line)
		if err != nil {
			fmt.Println("Failed to decode hex: ", err)
			return false
		}

		var score float64 = 0
		blocks := len(bytes) / 16
		shifted := bytes
		for i := 0; i < blocks; i++ {
			// Shift by one block.
			shifted = append(shifted[16:], shifted[:16]...)

			h, err := hammingDistance(shifted, bytes)
			if err != nil {
				fmt.Println("Failed to calculate hamming distance:", err)
				return false
			}

			score += float64(h)
		}
		score /= float64(blocks * blocks)

		//fmt.Println(lnum, "-", score)

		if score < bestScore {
			bestScore = score
			bestLine = lnum
		}
	}

	fmt.Println("AES_ECB on line:", bestLine)
	return true
}
