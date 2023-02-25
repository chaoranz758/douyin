package mysql

import "douyin/service/user/model"

func createTable() {
	_ = Db.AutoMigrate(&model.User{})
}
