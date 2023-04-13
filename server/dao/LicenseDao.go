/*
 * package dao 数据访问层
 */
package dao

import (
	"server/pkg/mysql"
	"time"
)

type License struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	License     string    `gorm:"column:license" json:"license"`
	Signature   string    `gorm:"column:signature" json:"signature"`
	Type        string    `gorm:"column:type" json:"type"`
	Object      string    `gorm:"column:object" json:"object"`
	Expiration  time.Time `gorm:"column:expiration" json:"expiration"`
	Project     string    `gorm:"column:project" json:"project"`
	Module      string    `gorm:"column:module" json:"module"`
	CreatedTime time.Time `gorm:"column:created_time" json:"created_time"`
	UpdatedTime time.Time `gorm:"column:updated_time" json:"updated_time"`
}

/*
 * SaveLicenseToDatabase 将许可证保存到数据库中
 *
 * @params: license *License - 许可证结构体指针
 * @returns: error - 任何可能发生的错误
 */
func AddLicenseToDatabase(license *License) error {

	db, _ := mysql.ConnectionMySQL()

	if err := db.Table("license_info").Create(&license).Error; err != nil {
		return err
	}

	return nil
}
