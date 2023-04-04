package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"server/config"
	"strconv"
)

func ConnectionMySQL() (*gorm.DB, error) {

	user := config.GetConfig("mysql", "user")
	password := config.GetConfig("mysql", "password")
	host := config.GetConfig("mysql", "host")
	port := config.GetConfig("mysql", "port")
	database := config.GetConfig("mysql", "database")
	charset := config.GetConfig("mysql", "charset")
	parseTimeStr := config.GetConfig("mysql", "parseTime")

	parseTime, err := strconv.ParseBool(parseTimeStr)
	if err != nil {
		parseTime = true // 默认值为true
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%v&loc=Local", user, password, host, port, database, charset, parseTime)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get DB object: %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
