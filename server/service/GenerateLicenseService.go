/*
 * Package service 提供生成许可证的服务
 * Info - 包含授权详细信息
 * Msg - 包含Info对象、状态和代码
 * License - 包含许可证ID、发放日期、类型、过期日期、允许用户数量、有效天数、项目名称和模块名称
 */

package service

import (
	"math/rand"
	"time"
)

type Info struct {
	Id            string `json:"id"`
	License       string `json:"license"`
	SignatureCode string `json:"signatureCode"`
	Date          string `json:"date"`
	Type          string `json:"type"`
	Expiration    string `json:"expiration"`
	AllowedUsers  string `json:"usersNum"`
	Project       string `json:"project"`
	Module        string `json:"module"`
}

type Msg struct {
	Info   Info   `json:"info"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type License struct {
	ID             string
	LicenseID      string
	Date           time.Time
	SignatureCode  string
	Type           string
	ExpirationDate time.Time
	AllowedUsers   uint
	Project        string
	Module         string
}

/*
 * GenerateLicense 生成许可证token
 *
 * @params: licenseType string - 许可证类型
 *			signatureCode string - 机器特征码
 * 			expiration time.Time - 过期日期
 *			allowedUsers uint - 允许的用户数量
 * 			days uint - 许可证有效天数
 *			obj string - 项目名称
 *			module string - 模块名称
 * @returns:License - 指向生成的license对象的指针
 * 			error - 任何可能发生的错误
 */
func GenerateLicense(signatureCode string, licenseType string, expiration time.Time, allowedUsers uint, obj string, module string) (*License, error) {

	// 生成随机、唯一的license
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "uBIXDoeA7GkpSKcNfuXbBLTs7nK8bJaNPEbN2zR5DLqs373P4uKWDxNuqfyYLY5XTWng4QB4"

	b := make([]byte, 72)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	licenseID := string(b)

	license := &License{
		LicenseID:      licenseID,
		Date:           time.Now().UTC(),
		SignatureCode:  signatureCode,
		Type:           licenseType,
		ExpirationDate: expiration,
		AllowedUsers:   allowedUsers,
		Project:        obj,
		Module:         module,
	}

	// 这里可以添加保存到数据库或文件等持久化操作
	// ...
	//err := ioutil.WriteFile("license.txt", []byte(license.LicenseID), 0644)
	//if err != nil {
	//	// 处理错误
	//	fmt.Println("Error:", err)
	//	return nil, nil
	//}
	return license, nil
}
