package main

// The main challenge runner.

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

var challenges = [][]func() bool{
	set1,
}

func main() {
	flag.Parse()

	for _, a := range flag.Args() {
		challenge, err := parseChallenge(a)
		if err != nil {
			fmt.Println("Invalid argument \"", a, "\": ", err)
			continue
		}

		fmt.Println("Running challenge ", a)
		if challenge() {
			fmt.Println("Challenge ", a, " passed!")
		} else {
			fmt.Println("Challenge ", a, " failed!")
		}
	}
}

func parseChallenge(s string) (func() bool, error) {
	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid format")
	}

	set, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return nil, errors.New("invalid format")
	}

	num, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return nil, errors.New("invalid format")
	}

	set--
	num--

	if int(set) >= len(challenges) {
		return nil, errors.New("invalid set")
	}

	if int(num) >= len(challenges[set]) {
		return nil, errors.New("invalid challenge in set")
	}

	return challenges[set][num], nil
}
