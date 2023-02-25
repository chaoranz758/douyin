package mysql

import (
	"douyin/service/comment/model"
	"errors"
	"go.uber.org/zap"
)

const (
	errorCreateComment0     = "评论数据没有插入进comment表但没有报错"
	errorCreateComment      = "create comment failed"
	errorDeleteCommentInfo0 = "评论数据没有从comment表删除但没有报错"
	errorDeleteCommentInfo  = "delete comment information failed"
)

var (
	errorCreateComment0V     = errors.New("评论数据没有插入进comment表但没有报错")
	errorDeleteCommentInfo0V = errors.New("评论数据没有从comment表删除但没有报错")
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
		zap.L().Error("revert comment failed not found", zap.Error(errors.New("revert comment failed not found")))
		return errors.New("revert comment failed not found")
	}
	if result.Error != nil {
		zap.L().Error("revert comment failed", zap.Error(result.Error))
		return result.Error
	}
	return nil
}
