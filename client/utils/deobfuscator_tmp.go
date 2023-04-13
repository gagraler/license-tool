package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"unicode"
)

func DeobfuscateXor(input string, signatureCode string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode input: %s", err)
	}

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	plaintext := make([]byte, len(ciphertext))

	for i := 0; i < len(ciphertext); i++ {
		plaintext[i] = ciphertext[i] ^ scBytes[i%len(scBytes)]
	}

	return string(plaintext), nil
}

func DeobfuscateRandomReplace(input string, signatureCode string, ratio float64) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode input: %s", err)
	}

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	plaintext := make([]byte, len(ciphertext))
	replaceTable := make(map[byte]byte)

	for i := 0; i < len(ciphertext); i++ {
		originalByte := ciphertext[i] ^ scBytes[i%len(scBytes)]

		if replacedByte, ok := replaceTable[originalByte]; ok {
			plaintext[i] = replacedByte
		} else {
			if rand.Float64() < ratio {
				randByte := byte(rand.Intn(256))
				for randByte == originalByte || randByte == '0' || randByte == '1' || unicode.IsLetter(rune(randByte)) {
					randByte = byte(rand.Intn(256))
				}
				plaintext[i] = randByte
				replaceTable[originalByte] = randByte
			} else {
				plaintext[i] = originalByte
			}
		}
	}

	return string(plaintext), nil
}

func DeobfuscateReverseXor(input string, signatureCode string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode input: %s", err)
	}

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	plaintext := make([]byte, len(ciphertext))

	for i := 0; i < len(ciphertext); i++ {
		plaintext[len(ciphertext)-1-i] = ciphertext[i] ^ scBytes[i%len(scBytes)]
	}

	return string(plaintext), nil
}

func DeobfuscateBase64Xor(input string, signatureCode string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return "", fmt.Errorf("failed to decode input: %s", err)
	}

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	plaintext := make([]byte, len(ciphertext))

	for i := 0; i < len(ciphertext); i++ {
		plaintext[i] = ciphertext[i] ^ scBytes[i%len(scBytes)]
	}

	decoded, err := base64.URLEncoding.DecodeString(string(plaintext))
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode: %s", err)
	}

	return string(decoded), nil
}
