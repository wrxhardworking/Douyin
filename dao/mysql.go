package dao

import (
	"douyin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 内涵连接池管理
var db *gorm.DB

// DbInit  数据库初始函数
func DbInit() error {
	var err error
	dsn := config.C.Mysql.Username + ":" + config.C.Mysql.Password + "@tcp(" + config.C.Mysql.Ipaddress + ":" + config.C.Mysql.Port + ")/" + config.C.Mysql.Dbname + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return err
}
