package mysql

import (
	"douyin/service/video/model"
	"go.uber.org/zap"
)

func CreateVideoInfo(v *model.Video) error {
	result := db.Create(v)
	if result.RowsAffected == 0 {
		zap.L().Error(errorCreateVideo0, zap.Error(errorCreateUser0V))
		return errorCreateUser0V
	}
	if result.Error != nil {
		zap.L().Error(errorCreateVideo, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func DeleteVideoInfo(videoId int64) error {
	result := db.Where("video_id = ?", videoId).Delete(&model.Video{})
	if result.RowsAffected == 0 {
		zap.L().Error(errorDeleteVideo0, zap.Error(errorDeleteUser0V))
		return errorDeleteUser0V
	}
	if result.Error != nil {
		zap.L().Error(errorDeleteVideo, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetVideoInfo(userID int64, vs *[]model.Video) error {
	result := db.Where("user_id=?", userID).Find(vs)
	if result.Error != nil {
		zap.L().Error(errorGetVideoInfo, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetVideoByVideoId(video *model.Video, videoId int64) error {
	return db.Where("video_id = ?", videoId).First(video).Error
}

func GetVideoListByVideoId(video *[]model.Video, videoId []int64) error {
	return db.Where("video_id in ?", videoId).Find(video).Error
}

func GetVideoListByCreateTime(video *[]model.Video, createTime int64, count int) error {
	return db.Where("create_time < ?", createTime).Limit(count).Order("create_time desc").Find(video).Error
}
