package main

import (
	"crypto/aes"
	"errors"
)

// Decrypt AES in ECB mode.
func decryptAesEcb(cipherbytes, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	bs := c.BlockSize()
	if len(cipherbytes)%bs != 0 {
		return nil, errors.New("length of input must be multiple of blocksize.")
	}

	plainbytes := make([]byte, len(cipherbytes))
	p := plainbytes[:]
	for len(cipherbytes) > 0 {
		c.Decrypt(p, cipherbytes)
		p = p[bs:]
		cipherbytes = cipherbytes[bs:]
	}

	return plainbytes, nil
}
