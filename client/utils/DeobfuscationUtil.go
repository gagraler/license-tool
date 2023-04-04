/*
 * package utils 工具包，提供了一些常用的工具函数
 */
package utils

import (
	"encoding/base64"
	"errors"
)

/*
 * DeobfuscationUtil 字节数组反混淆工具函数，用于将输入的字节数组与特征码进行异或运算并进行反加混淆
 * @param: input []byte - 需要加密的字节数组
 *		   signatureCode string - 签名代码
 * @return: string - 加密后的字符串
 * 		 	error - 如果签名代码为空，则返回一个错误
 */
func DeobfuscationUtil(obfuscatedStr string, signatureCode string) ([]byte, error) {

	if signatureCode == "" {
		return nil, errors.New("signature code is empty")
	}

	scBytes := []byte(signatureCode)

	// 对传入的被混淆字符串进行 base64 解码。
	obfBytes, err := base64.URLEncoding.DecodeString(obfuscatedStr)
	if err != nil {
		return nil, err
	}

	// 创建一个与解码后的字节数组长度相同的字节数组，用于存储还原后的数据。
	output := make([]byte, len(obfBytes))

	// 对每个字节进行按位异或运算（XOR operation），以还原出原始的数据。
	for i := 0; i < len(obfBytes); i++ {
		output[i] = obfBytes[i] ^ scBytes[i%len(scBytes)]
	}

	// 返回还原后的字节数组和一个可能存在的错误。
	return output, nil
}
