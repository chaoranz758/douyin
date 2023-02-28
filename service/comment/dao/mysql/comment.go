package mysql

import (
	"douyin/service/comment/model"
	"go.uber.org/zap"
)

func PushCommentInfo(comment model.Commit) error {
	result := db.Create(&comment)
	if result.RowsAffected == 0 {
		zap.L().Error(errorCreateComment0, zap.Error(errorCreateComment0V))
		return errorCreateComment0V
	}
	if result.Error != nil {
		zap.L().Error(errorCreateComment, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func DeleteCommentInfo(commentId int64) error {
	result := db.Where("commit_id = ?", commentId).Delete(&model.Commit{})
	if result.RowsAffected == 0 {
		zap.L().Error(errorDeleteCommentInfo0, zap.Error(errorDeleteCommentInfo0V))
		return errorDeleteCommentInfo0V
	}
	if result.Error != nil {
		zap.L().Error(errorDeleteCommentInfo, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetCommentInfo(cs *[]model.Commit, videoId int64) error {
	return db.Where("video_id = ?", videoId).Find(cs).Error
}

func RevertComment(commentId int64) error {
	result := db.Where("commit_id = ?", commentId).Update("deleted_at", nil)
	if result.RowsAffected == 0 {
		zap.L().Error(errorRevertCommentInfo0, zap.Error(errorRevertCommentInfo0V))
		return errorRevertCommentInfo0V
	}
	if result.Error != nil {
		zap.L().Error(errorRevertCommentInfo, zap.Error(result.Error))
		return result.Error
	}
	return nil
}
