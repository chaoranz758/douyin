package mysql

import "errors"

const (
	errorFavoriteRelation0       = "点赞关系数据没有插入进favorite表但没有报错"
	errorFavoriteRelation        = "create favorite relation failed"
	errorDeleteFavoriteRelation0 = "点赞关系没有在favorite表中删除成功"
	errorDeleteFavoriteRelation  = "delete favorite relation from favorite table failed"
	errorGetUserFavoriteID       = "get user favorite id failed"
	errorConnectDB               = "connect DB failed"
	errorOpenDB                  = "open db.DB() failed"
)

var (
	errorFavoriteRelation0V       = errors.New("点赞关系数据没有插入进favorite表但没有报错")
	errorDeleteFavoriteRelation0V = errors.New("点赞关系没有在favorite表中删除成功")
)
