/*
 * package service 业务函数包，处理业务逻辑
 */
package service

import (
	"log"
	"math/rand"
	"server/dao"
	"server/pkg/mysql"
	"server/util"
	"time"
)

/*
 * GenerateLicense 生成许可证token
 *
 * @params: params 请求参数结构体
 * @returns: License - 指向生成的license对象的指针
 * 			 error - 任何可能发生的错误
 */
func GenerateLicenseService(params *dao.License) (*dao.License, error) {

	// 生成随机、唯一的license
	rand.Seed(time.Now().UnixNano())

	const letterBytes = "uBIXDoeA7GkpSKcNfuXbBLTs7nK8bJaNPEbN2zR5DLqs373P4uKWDxNuqfyYLY5XTWng4QB4"

	b := make([]byte, 72)

	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	licenseStr := string(b)

	// 构建license消息体
	license := &dao.License{
		ID:         util.GenerateUniqueID(),
		License:    licenseStr,
		Signature:  params.Signature,
		Object:     params.Object,
		Type:       params.Type,
		Expiration: params.Expiration,
		Project:    params.Project,
		Module:     params.Module,
		CreatedAt:  time.Now().UTC(),
	}

	// 保存到数据库
	_, err := mysql.ConnectionMySQL()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := dao.AddLicenseToDatabase(license); err != nil {
		return nil, err
	}

	return license, nil
}
