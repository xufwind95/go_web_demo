package model

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // 下划线为仅仅调用该包的Init函数
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Database struct {
	Self   *gorm.DB
	Docker *gorm.DB
}

// 定义数据库操作对象
var DB *Database

// 连接数据库
func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		"Local",
	)
	log.Info(config)

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database connection failed.Database name: %s", name)
	}

	setupDB(db)

	return db
}

// 数据库相关设置
func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	db.DB().SetMaxIdleConns(0)
}

// 初始化 self 数据库
func InitSelfDB() *gorm.DB {
	return openDB(
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"),
	)
}

// 获取 self 数据库操作对象
func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

// 初始化 docker 数据库
func InitDockerDB() *gorm.DB {
	return openDB(
		viper.GetString("docker_db.username"),
		viper.GetString("docker_db.password"),
		viper.GetString("docker_db.addr"),
		viper.GetString("docker_db.name"),
	)
}

// 获取 docker 数据库操作对象
func GetDockerDB() *gorm.DB {
	return InitDockerDB()
}

// 数据库初始化入口
func (db *Database) Init() {
	DB = &Database{
		Self:   GetSelfDB(),
		Docker: GetDockerDB(),
	}
}

// 关闭数据库连接
func (db *Database) Close() {
	DB.Self.Close()
	DB.Docker.Close()
}
