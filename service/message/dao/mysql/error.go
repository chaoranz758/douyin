package mysql

import "errors"

const (
	errorCreateMessage0 = "消息没有插入进message表但没有报错"
	errorCreateMessage  = "create message failed"
	errorConnectDB      = "connect DB failed"
	errorOpenDB         = "open db.DB() failed"
	errorDeleteMessage0 = "delete message not found"
	errorDeleteMessage  = "delete message failed"
)

var (
	errorCreateMessage0V = errors.New("消息没有插入进message表但没有报错")
	errorDeleteMessage0V = errors.New("delete message not found")
)
