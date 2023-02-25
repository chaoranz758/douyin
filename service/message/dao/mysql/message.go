package mysql

import (
	"douyin/service/message/model"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	errorCreateMessage0 = "消息没有插入进message表但没有报错"
	errorCreateMessage  = "create message failed"
)

var (
	errorCreateMessage0V = errors.New("消息没有插入进message表但没有报错")
)

func CreateMessage(m *model.Message) error {
	result := db.Create(m)
	if result.RowsAffected == 0 {
		zap.L().Error(errorCreateMessage0, zap.Error(errorCreateMessage0V))
		return errorCreateMessage0V
	}
	if result.Error != nil {
		zap.L().Error(errorCreateMessage, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func DeleteMessage(messageId int64) error {
	result := db.Where("message_id = ?", messageId).Delete(&model.Messages{})
	if result.RowsAffected == 0 {
		zap.L().Error("delete message not found", zap.Error(errors.New("delete message not found")))
		return errors.New("delete message not found")
	}
	if result.Error != nil {
		zap.L().Error("delete message failed", zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetMessage(ms *[]model.Message, userId, toUserId, t int64, isFirst bool) error {
	if isFirst == true {
		return db.Where("user_id = ? and to_user_id = ?", userId, toUserId).Find(ms).Error
	}
	return db.Where("user_id = ? and to_user_id = ? and created_at > ?", userId, toUserId, t).Find(ms).Error
}

func GetMessageByUserId(ms *[]model.Message, userId int64) error {
	return db.Where("user_id = ?", userId).Find(ms).Error
}

func GetMessageLast(ms *[]model.Message, userId, toUserId int64) error {
	if err := db.Where("user_id = ? and to_user_id = ?", userId, toUserId).Last(ms).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			zap.L().Info(fmt.Sprintf("用户%v与用户%v之间没有聊过天", userId, toUserId))
			return nil
		}
		return err
	}
	return nil
}
