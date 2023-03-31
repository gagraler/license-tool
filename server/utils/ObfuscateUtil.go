package utils

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"unicode"
)

func ObfuscationUtil(input []byte, signatureCode string) (string, error) {

	const targetLen = 4096

	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, targetLen)
	for i := 0; i < targetLen; i++ {
		if i < len(input) {
			ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
		} else {
			randByte := byte(rand.Intn(256))
			for randByte != '0' && randByte != '1' && !unicode.IsLetter(rune(randByte)) {
				randByte = byte(rand.Intn(256))
			}
		}
	}

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

//func ObfuscationUtil(input []byte, signatureCode string) (string, error) {
//	const targetLen = 4096
//
//	if signatureCode == "" {
//		return "", errors.New("signature code is empty")
//	}
//
//	scBytes := []byte(signatureCode)
//	ciphertext := make([]byte, targetLen)
//	for i := 0; i < targetLen; i++ {
//		if i < len(input) {
//			ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
//		} else {
//			ciphertext[i] = byte(rand.Intn(256))
//		}
//	}
//
//	return base64.URLEncoding.EncodeToString(ciphertext), nil
//}
