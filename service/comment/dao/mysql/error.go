package mysql

import "errors"

const (
	errorCreateComment0     = "评论数据没有插入进comment表但没有报错"
	errorCreateComment      = "create comment failed"
	errorDeleteCommentInfo0 = "评论数据没有从comment表删除但没有报错"
	errorDeleteCommentInfo  = "delete comment information failed"
	errorRevertCommentInfo0 = "revert comment failed not found"
	errorRevertCommentInfo  = "revert comment failed"
	errorConnectDB          = "connect DB failed"
	errorOpenDB             = "open db.DB() failed"
)

var (
	errorCreateComment0V     = errors.New("评论数据没有插入进comment表但没有报错")
	errorDeleteCommentInfo0V = errors.New("评论数据没有从comment表删除但没有报错")
	errorRevertCommentInfo0V = errors.New("revert comment failed not found")
)
