// @author Zongsheng Xu 2023/3/31 1:18

package utils

import (
	"encoding/base64"
	"errors"
)

func DeobfuscateUtil(ciphertext string, signatureCode string) ([]byte, error) {
	if signatureCode == "" {
		return nil, errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)

	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	plaintext := make([]byte, len(decodedCiphertext))
	for i := 0; i < len(decodedCiphertext); i++ {
		plaintext[i] = decodedCiphertext[i] ^ scBytes[i%len(scBytes)]
	}

	return plaintext, nil
}
