package mysql

import (
	"douyin/proto/user/request"
	"douyin/service/user/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

func UserRegister(u *model.User) error {
	result := Db.Create(u)
	if result.RowsAffected == 0 {
		zap.L().Error(errorCreateUser0, zap.Error(errorCreateUser0V))
		return errorCreateUser0V
	}
	if result.Error != nil {
		zap.L().Error(errorCreateUser, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func UserLogin(u *model.User, req *request.DouyinUserLoginRequest) error {
	if err := Db.Where("user_name = ?", req.Username).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.L().Error(errorUserLogin0, zap.Error(errorUserLogin0V))
			return errorUserLogin0V
		}
		zap.L().Error(errorUserLogin, zap.Error(err))
		return err
	}
	return nil
}

func GetUserInfo(userID int64, u *model.User) error {
	if err := Db.Where("user_id = ?", userID).First(u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.L().Error(errorUserInfoNotExist, zap.Error(err))
			return err
		}
		zap.L().Error(errorGetUserInfoFailed, zap.Error(err))
		return err
	}
	return nil
}

func GetUserInfoList(userID []int64, u *[]model.User) error {
	var users []string
	for i := 0; i < len(userID); i++ {
		s := strconv.FormatInt(userID[i], 10)
		users = append(users, s)
	}
	if err := Db.Where("user_id IN ?", users).Find(u).Error; err != nil {
		zap.L().Error(errorGetUserInfoListFailed, zap.Error(err))
		return err
	}
	return nil
}
