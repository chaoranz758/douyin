package mysql

import "douyin/service/follow/model"

func createTable() {
	_ = db.AutoMigrate(&model.Follow{})
}
