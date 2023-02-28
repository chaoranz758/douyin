package mysql

import "errors"

const (
	errorFollowRelation0 = "关注关系数据没有插入进follow表但没有报错"
	errorFollowRelation  = "create follow relation failed"
	errorDeleteFollow0   = "关注关系在follow表中没有被删除成功"
	errorDeleteFollow    = "delete follow failed"
	errorConnectDB       = "connect DB failed"
	errorOpenDB          = "open db.DB() failed"
)

var (
	errorFollowRelation0V = errors.New("关注关系数据没有插入进follow表但没有报错")
	errorDeleteFollow0V   = errors.New("关注关系在follow表中没有被删除成功")
)
