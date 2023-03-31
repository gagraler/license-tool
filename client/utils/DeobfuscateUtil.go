package utils

import (
	"encoding/base64"
	"errors"
)

/* DeobfuscateUtil 对密文进行逆向操作，得到原始明文数据
 * @params: ciphertext: 待解密的密文字符串
 * 			signatureCode: 用于加密的签名代码，必须非空
 * @return: []byte: 如果解密成功，则返回一个字节数组表示的明文数据；否则为 nil
 * 			error: 如果解密失败，则返回一个错误对象；否则为 nil
 */
func DeobfuscateUtil(ciphertext string, signatureCode string) ([]byte, error) {

	if signatureCode == "" {
		return nil, errors.New("signature code is empty")
	}

	// 将 signatureCode 转换为字节数组 scBytes
	scBytes := []byte(signatureCode)

	// 使用 base64 解码密文数据
	decodedCiphertext, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	// 对密文进行逆向操作，得到原始明文数据
	plaintext := make([]byte, len(decodedCiphertext))
	for i := 0; i < len(decodedCiphertext); i++ {
		plaintext[i] = decodedCiphertext[i] ^ scBytes[i%len(scBytes)]
	}

	return plaintext, nil
}
