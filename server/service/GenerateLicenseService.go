/*
 * package service 业务函数包，处理业务逻辑
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
	Object        string `json:"object"`
	Project       string `json:"project"`
	Module        string `json:"module"`
}

type Msg struct {
	Info   Info   `json:"info"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type LicenseParams struct {
	ID            string
	LicenseID     string
	Date          time.Time
	SignatureCode string
	Type          string
	Object        string
	Expiration    time.Time
	Project       string
	Module        string
}

/*
 * GenerateLicense 生成许可证token
 *
 * @params: licenseType string - 许可证类型
 *			signatureCode string - 机器特征码
 * 			expiration time.Time - 过期日期
 *			Object uint - 允许的用户数量
 * 			days uint - 许可证有效天数
 *			obj string - 项目名称
 *			module string - 模块名称
 * @returns:License - 指向生成的license对象的指针
 * 			error - 任何可能发生的错误
 */
func GenerateLicense(signatureCode string, licenseType string, expiration time.Time, object string, project string, module string) (*LicenseParams, error) {

	// 生成随机、唯一的license
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "uBIXDoeA7GkpSKcNfuXbBLTs7nK8bJaNPEbN2zR5DLqs373P4uKWDxNuqfyYLY5XTWng4QB4"

	b := make([]byte, 72)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	licenseID := string(b)

	license := &LicenseParams{
		LicenseID:     licenseID,
		Date:          time.Now().UTC(),
		SignatureCode: signatureCode,
		Object:        object,
		Type:          licenseType,
		Expiration:    expiration,
		Project:       project,
		Module:        module,
	}

	// 这里可以添加保存到数据库或文件等持久化操作

	return license, nil
}
