package config

import (
	"fmt"
	"net/url"
	"servermanager/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%s&loc=%s",
		Config.Database.MySQL.User,
		Config.Database.MySQL.Password,
		Config.Database.MySQL.Host,
		Config.Database.MySQL.Port,
		Config.Database.MySQL.DBName,
		Config.Database.MySQL.Charset,
		fmt.Sprintf("%t", Config.Database.MySQL.ParseTime),
		url.QueryEscape(Config.Database.MySQL.Loc),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 自动迁移数据库结构
	err = db.AutoMigrate(&models.Material{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
