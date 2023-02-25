package mysql

import "douyin/service/favorite/model"

func createTable() {
	_ = db.AutoMigrate(&model.Favorite{})
}
