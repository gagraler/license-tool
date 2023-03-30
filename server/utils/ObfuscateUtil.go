// Author: Zongsheng Xu 2023/3/30 21:34

package utils

import (
	"encoding/base64"
	"errors"
)

func ObfuscateUtil(input []byte, signatureCode string) (string, error) {
	if signatureCode == "" {
		return "", errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)
	ciphertext := make([]byte, len(input))
	for i := 0; i < len(input); i++ {
		ciphertext[i] = input[i] ^ scBytes[i%len(scBytes)]
	}
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}
