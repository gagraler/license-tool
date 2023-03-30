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
	License          string `json:"license"`
	Date             string `json:"date"`
	AuthorizedObject string `json:"authorized_object"`
	Expiration       string `json:"expiration"`
	AllowedUsers     string `json:"allowed_users"`
	Days             string `json:"days"`
	Project          string `json:"project"`
	Module           string `json:"module"`
}

type Msg struct {
	Info   Info   `json:"info"`
	Status string `json:"status"`
	Code   int    `json:"code"`
}

type License struct {
	LicenseID      string
	Date           time.Time
	Type           string
	ExpirationDate time.Time
	AllowedUsers   uint
	Days           uint
	Project        string
	Module         string
}

/*
 * GenerateLicense 生成许可证token
 *
 * @params: licenseType string - 许可证类型
 * 			expiration time.Time - 过期日期
 *			allowedUsers uint - 允许的用户数量
 * 			days uint - 许可证有效天数
 *			obj string - 项目名称
 *			module string - 模块名称
 * @returns:License - 指向生成的license对象的指针
 * 			error - 任何可能发生的错误
 */
func GenerateLicense(licenseType string, expiration time.Time, allowedUsers uint, days uint, obj string, module string) (*License, error) {

	// 生成随机、唯一的license
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "TS8dqQifBcZxBZZRdQQaflbKqvkqGT1KdevohEOtutYhftB5PIvcsDIcUbY3XWDuapAD1al8e9dNBqL3IuMN2gvCzHX7hgs8VosnpWltRPdo3BGbpol42muV8LTREpjjnyN1uEffK1HO8P6WKmXkGoWUOAvdfe0zFedPHJVEp981TCZJmWhs65N4uYKD7Zv1HlPtHmLe7T2qt4UyYC7Lx65q7mDVjnTArVPhz7k39WL7nx0cDovNCjQwRmTXqSOg"

	b := make([]byte, 256)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	licenseID := string(b)

	license := &License{
		LicenseID:      licenseID,
		Date:           time.Now().UTC(),
		Type:           licenseType,
		ExpirationDate: expiration,
		AllowedUsers:   allowedUsers,
		Days:           days,
		Project:        obj,
		Module:         module,
	}

	// 这里可以添加保存到数据库或文件等持久化操作
	// ...

	return license, nil
}
