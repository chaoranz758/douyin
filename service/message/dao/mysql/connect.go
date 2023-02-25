package mysql

import (
	"douyin/service/message/config"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func Init() (err error) {
	//连接到mysql数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.Mysql.User,
		config.Config.Mysql.Password,
		config.Config.Mysql.Host,
		config.Config.Mysql.Port,
		config.Config.Mysql.Dbname,
	)
	println()
	//Logger: logger.Default.LogMode(logger.Info
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("open db.DB() failed", zap.Error(err))
		return err
	}
	//空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.Config.Mysql.MaxIdleConns)
	//打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(config.Config.Mysql.MaxOpenConns)
	createTable()
	return nil
}
