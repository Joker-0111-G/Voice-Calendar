package models

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:G20050111g@tcp(127.0.0.1:3306)/voice_calendar?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("MySQL 连接失败: %v", err)
	}

	log.Println("MySQL 连接成功，执行自动建表...")
	err = DB.AutoMigrate(&Event{})
	if err != nil {
		log.Fatalf("表结构迁移失败: %v", err)
	}
}