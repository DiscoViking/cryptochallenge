package main

import "math"

// Frequency Analysis Function

// English letter frequencies.
var letterFrequencies = map[byte]float64{
	'a': 0.08167,
	'b': 0.01492,
	'c': 0.02782,
	'd': 0.04253,
	'e': 0.12702,
	'f': 0.02228,
	'g': 0.02015,
	'h': 0.06094,
	'i': 0.06966,
	'j': 0.00153,
	'k': 0.00772,
	'l': 0.04025,
	'm': 0.02406,
	'n': 0.06749,
	'o': 0.07507,
	'p': 0.01929,
	'q': 0.00095,
	'r': 0.05987,
	's': 0.06327,
	't': 0.09056,
	'u': 0.02758,
	'v': 0.00978,
	'w': 0.02360,
	'x': 0.00150,
	'y': 0.01974,
	'z': 0.00074,
}

var knownPunctuation = map[byte]struct{}{
	'.':  struct{}{},
	':':  struct{}{},
	';':  struct{}{},
	'?':  struct{}{},
	'/':  struct{}{},
	'\'': struct{}{},
	'"':  struct{}{},
	'!':  struct{}{},
	'@':  struct{}{},
	'#':  struct{}{},
}

// Calculates byte frequences for a bytestring.
func byteFrequencies(in []byte) map[byte]float64 {
	freqs := map[byte]float64{}
	for _, b := range in {
		if _, ok := freqs[b]; ok {
			freqs[b]++
		} else {
			freqs[b] = 1
		}
	}

	for k, v := range freqs {
		freqs[k] = v / float64(len(in))
	}

	return freqs
}

// Scores a bytestring for its similarity to English text.
// Uses total square distance for each letter frequency.
func englishness(in []byte) float64 {
	var score float64 = 0

	freqs := byteFrequencies(in)

	// Add in all the english characters so we definitely count them.
	for c, _ := range letterFrequencies {
		if _, ok := freqs[c]; !ok {
			freqs[c] = 0
		}
	}

	for c, f := range freqs {
		x, ok := letterFrequencies[c]
		if !ok {
			x = 0
		}
		score += math.Sqrt(f * x)
	}

	return score
}
