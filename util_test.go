package main

// Random Tests

import (
	"testing"
)

func TestHamming(t *testing.T) {
	a := []byte("this is a test")
	b := []byte("wokka wokka!!!")
	d, err := hammingDistance(a, b)
	if err != nil {
		t.Fatal(err.Error())
	}

	if d != 37 {
		t.Fatalf("Expected 37, got %v\n", d)
	}
}
