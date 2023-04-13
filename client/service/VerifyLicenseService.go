/*
 * package service 业务函数包，处理业务逻辑
 */
package service

import (
	"client/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

/*
 * ReadLicenseFile 函数用于读取许可证文件的内容。
 * @params: licenseName string：许可证文件名
 * @return: []byte：许可证文件的内容，以 byte 数组的形式返回
 * 			error：如果读取文件出错，则返回一个错误对象；否则为 nil
 */
func ReadLicenseFile(licenseName string) ([]byte, error) {

	// 打开文件
	ciphertext, err := os.Open(licenseName)
	if err != nil {
		return nil, err
	}

	// 延迟关闭文件
	defer func(ciphertext *os.File) {
		err := ciphertext.Close()
		if err != nil {
			fmt.Println("failed to close file:", err)
		}
	}(ciphertext)

	// 读取文件内容
	licenseContent, err := ioutil.ReadAll(ciphertext)
	if err != nil {
		return nil, err
	}

	return licenseContent, nil
}

/*
 * DeobfuscateLicense 用于对许可证文件内容进行反混淆和提取。
 * @params: licenseContent：需要反混淆和提取的许可证文件内容。
 * @return: 如果成功反混淆和提取许可证文件内容，则返回包含 JSON 字符串的字符串和 nil；否则返回空字符串和错误信息
 */
func DeobfuscateLicense(licenseContent []byte) (string, error) {

	// 调用反混淆工具函数对许可内容反混淆
	outputBytes, err := utils.DeobfuscationUtil(string(licenseContent), utils.MachineCode())
	if err != nil {
		fmt.Println("Error obfuscating data:", err)
		return "", err
	}

	// 提取json字符串
	re := regexp.MustCompile(`{(.*)}`)
	jsonString := re.FindStringSubmatch(string(outputBytes))[0]
	if err != nil {
		return "", err
	}

	fmt.Println("Obfuscating bytes:", jsonString)

	// 返回包含 JSON 数据的字符串和 nil
	return jsonString, nil
}

/*
 * ExtractSignatureCode 用于从 JSON 字符串中提取 signatureCode
 * @params: jsonString：包含 JSON 字符串的字符串
 * @return: 如果成功提取 signatureCode，则返回 signatureCode 和 nil；否则返回空字符串和错误信息
 */
func ExtractSignatureCode(jsonString string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return "", err
	}

	signatureCode, ok := data["signature"].(string)
	if !ok {
		return "", errors.New("failed to extract signatureCode from license")
	}

	return signatureCode, nil
}

/* VerifyLicense 验证许可文件
 * @params: licenseName: 许可文件名
 * @return: true 表示许可文件有效，false 表示许可文件无效
 *			error: 验证失败，则返回一个错误对象；否则为 nil
 */
func VerifyLicense(licenseName string) (bool, error) {
	// 读取许可文件内容
	licenseContent, err := ReadLicenseFile(licenseName)
	if err != nil {
		return false, err
	}

	// 反混淆许可文件
	jsonString, err := DeobfuscateLicense(licenseContent)
	if err != nil {
		return false, err
	}

	// 提取 signatureCode 值
	signatureCode, err := ExtractSignatureCode(jsonString)
	if err != nil {
		return false, err
	}

	// 校验 signatureCode 和 MachineCode 是否匹配
	machineCode := utils.MachineCode()
	if machineCode != signatureCode {
		return false, errors.New("license verification failed: signatureCode does not match MachineCode")
	} else {
		fmt.Println("license verification succeeded!")
		return true, nil
	}
}
