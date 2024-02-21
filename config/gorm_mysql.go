package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GORM_DB *gorm.DB

type Mysql struct {
	Username string `json:"username" gorm:"username"`
	Password string `json:"password" gorm:"password"`
	Host     string `json:"host" gorm:"host"`
	Port     string `json:"port" gorm:"port"`
}

func (m *Mysql) Dsn() string {
	return CONFIG.Mysql.Username + ":" + CONFIG.Mysql.Password + "@tcp(" + CONFIG.Mysql.Host + ":" + CONFIG.Mysql.Port + ")/msg_push_center?charset=utf8mb4&parseTime=True&loc=Local"
}

func GormMysql() *gorm.DB {
	var configMysql Mysql
	dsn := configMysql.Dsn()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}
