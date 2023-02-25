package mysql

import "douyin/service/video/model"

func createTable() {
	_ = db.AutoMigrate(&model.Video{})
}
