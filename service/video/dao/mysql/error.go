package mysql

import "errors"

const (
	errorCreateVideo0 = "视频数据没有插入进video表但没有报错"
	errorCreateVideo  = "create video failed"
	errorDeleteVideo0 = "视频数据没有从video表删除但没有报错"
	errorDeleteVideo  = "delete video failed"
	errorGetVideoInfo = "get video information failed"
	errorConnectDB    = "connect DB failed"
	errorOpenDB       = "open db.DB() failed"
)

var (
	errorCreateUser0V = errors.New("视频数据没有插入进video表但没有报错")
	errorDeleteUser0V = errors.New("视频数据没有从video表删除但没有报错")
)
