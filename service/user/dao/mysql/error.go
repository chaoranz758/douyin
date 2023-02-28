package mysql

import "errors"

const (
	errorCreateUser0           = "用户数据没有插入进user表但没有报错"
	errorCreateUser            = "create user failed"
	errorUserLogin0            = "在用户表中查不到对应的用户名"
	errorUserLogin             = "user login failed"
	errorUserInfoNotExist      = "user information not exist"
	errorGetUserInfoFailed     = "get user information failed"
	errorGetUserInfoListFailed = "get user information list failed"
	errorConnectDB             = "connect DB failed"
	errorOpenDB                = "open db.DB() failed"
)

var (
	errorCreateUser0V = errors.New("用户数据没有插入进user表但没有报错")
	errorUserLogin0V  = errors.New("在用户表中查不到对应的用户名")
)
