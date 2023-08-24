package database

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db   *gorm.DB
	once sync.Once
)

func init() {
	Init()
}

func Init() {
	once.Do(func() {
		initDb()
	})
}

func initDb() {
	// 配置信息
	host := "localhost"
	port := "3306"
	username := "root"
	password := "88888888"
	dbname := "douyin"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, ConnectErr := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if ConnectErr != nil {
		log.Println("数据库连接失败", ConnectErr)
		return
	}

	// 设置连接池
	sqlDB, GetErr := db.DB()
	if GetErr != nil {
		log.Println("获取底层数据库对象失败", GetErr)
		return
	}
	sqlDB.SetMaxIdleConns(10)   // 设置空闲连接池中的最大连接数
	sqlDB.SetMaxOpenConns(100)  // 设置打开数据库连接的最大数量
	sqlDB.SetConnMaxLifetime(0) // 连接可复用的最大时间

	Db = db
}
