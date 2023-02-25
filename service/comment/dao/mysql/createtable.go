package mysql

import "douyin/service/comment/model"

func createTable() {
	_ = db.AutoMigrate(&model.Commit{})
}
