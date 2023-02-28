package mysql

import (
	"douyin/service/message/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
		zap.L().Error(errorDeleteMessage0, zap.Error(errorDeleteMessage0V))
		return errorDeleteMessage0V
	}
	if result.Error != nil {
		zap.L().Error(errorDeleteMessage, zap.Error(result.Error))
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
			//zap.L().Info(fmt.Sprintf("用户%v与用户%v之间没有聊过天", userId, toUserId))
			return nil
		}
		return err
	}
	return nil
}
