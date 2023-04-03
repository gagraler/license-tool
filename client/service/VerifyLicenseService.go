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

/* VerifyLicense 验证许可文件
 * @params: licenseName: 许可文件名
 * @return: true 表示许可文件有效，false 表示许可文件无效
 *			error: 验证失败，则返回一个错误对象；否则为 nil
 */
func VerifyLicense(licenseName string) (bool, error) {

	var jsonString string

	// 打开许可文件
	ciphertext, err := os.Open(licenseName)
	if err != nil {
		panic(err)
	}

	// 异常处理
	defer func(ciphertext *os.File) {
		err := ciphertext.Close()
		if err != nil {
		}
	}(ciphertext)

	// 读取许可文件内容
	licenseContent, err := ioutil.ReadAll(ciphertext)
	if err != nil {
		panic(err)
	}

	outputBytes, err := utils.DeobfuscationUtil(string(licenseContent), utils.MachineCode())
	if err != nil {
		fmt.Println("Error deobfuscating data:", err)
		return false, nil
	} else {
		re := regexp.MustCompile(`{(.*)}`)
		jsonString = re.FindStringSubmatch(string(outputBytes))[0]
		if err != nil {
			panic(err)
		}
		fmt.Println("Deobfuscated bytes:", jsonString)
	}

	// 提取 signatureCode 值
	var (
		data map[string]interface{}
	)

	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return false, err
	}

	// 将 data["authorized"] 转换为 map[string]interface{} 类型，并赋值给 authorized 变量。
	// ok 的值为 true 表示转换成功，否则表示转换失败，需要返回错误信息。
	authorized, ok := data["authorized"].(map[string]interface{})
	if !ok {
		return false, errors.New("failed to extract authorized object from license")
	}

	// 将 authorized["signatureCode"] 转换为 string 类型，并赋值给 signatureCode 变量。
	// ok 的值为 true 表示转换成功，否则表示转换失败，需要返回错误信息。
	signatureCode, ok := authorized["signatureCode"].(string)
	if !ok {
		return false, errors.New("failed to extract signatureCode from license")
	}

	// 检查 signatureCode 和 MachineCode 是否匹配
	machineCode := utils.MachineCode()
	if machineCode != signatureCode {
		return false, errors.New("license verification failed: signatureCode does not match MachineCode")
	} else {
		fmt.Println("license verification succeeded!")
		return true, nil
	}
}
