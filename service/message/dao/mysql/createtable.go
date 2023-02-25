package mysql

import "douyin/service/message/model"

func createTable() {
	_ = db.AutoMigrate(&model.Message{})
}
