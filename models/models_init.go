package models

import "gorm.io/gorm"

// DB 全局数据库实例
var DB *gorm.DB

// SetDB 设置全局数据库实例
func SetDB(db *gorm.DB) {
	DB = db
}
