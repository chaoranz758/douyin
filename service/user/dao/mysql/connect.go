package mysql

import (
	"douyin/service/user/initialize/config"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
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
	//Logger: logger.Default.LogMode(logger.Info
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Error(errorConnectDB, zap.Error(err))
		return err
	}
	sqlDB, err := Db.DB()
	if err != nil {
		zap.L().Error(errorOpenDB, zap.Error(err))
		return err
	}
	//空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(config.Config.Mysql.MaxIdleConns)
	//打开数据库连接的最大数量
	sqlDB.SetMaxOpenConns(config.Config.Mysql.MaxOpenConns)
	//createTable()
	return nil
}
