package mysql

import (
	"douyin/service/favorite/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	favorite = iota
	cancelFavorite
)

func CreateFavoriteRelation(f *model.Favorite) error {
	var fo model.Favorite
	resultGet := db.Where("video_id = ? and user_id = ? and status = ?", f.VideoID, f.UserID, cancelFavorite).First(&fo)
	if resultGet.RowsAffected == 0 {
		result := db.Create(f)
		if result.RowsAffected == 0 {
			zap.L().Error(errorFavoriteRelation0, zap.Error(errorFavoriteRelation0V))
			return errorFavoriteRelation0V
		}
		if result.Error != nil {
			zap.L().Error(errorFavoriteRelation, zap.Error(result.Error))
			return result.Error
		}
		return nil
	}
	if resultGet.Error != nil {
		zap.L().Error(errorFavoriteRelation, zap.Error(resultGet.Error))
		return resultGet.Error
	}
	result := db.Model(&model.Favorite{}).Where("video_id = ? and user_id = ? and status = ?", f.VideoID, f.UserID, cancelFavorite).Update("status", favorite)
	if result.RowsAffected == 0 {
		zap.L().Error(errorFavoriteRelation0, zap.Error(errorFavoriteRelation0V))
		return errorFavoriteRelation0V
	}
	if result.Error != nil {
		zap.L().Error(errorFavoriteRelation, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func DeleteFavoriteRelation(videoID, userID int64) error {
	result := db.Model(&model.Favorite{}).Where("video_id = ? and user_id = ? and status = ?", videoID, userID, favorite).Update("status", cancelFavorite)
	if result.RowsAffected == 0 {
		zap.L().Error(errorDeleteFavoriteRelation0, zap.Error(errorDeleteFavoriteRelation0V))
		return errorDeleteFavoriteRelation0V
	}
	if result.Error != nil {
		zap.L().Error(errorDeleteFavoriteRelation, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetUserFavoriteID(userID int64, f *[]model.Favorite) error {
	result := db.Where("user_id = ? and status = ?", userID, favorite).Find(f)
	if result.Error != nil {
		zap.L().Error(errorGetUserFavoriteID, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetUserFavoriteBool(userId int64, videoIdList []int64) ([]bool, error) {
	var bs []bool
	println(len(videoIdList))
	for i := 0; i < len(videoIdList); i++ {
		var v model.Favorite
		if err := db.Where("video_id = ? and user_id = ? and status = ?", videoIdList[i], userId, favorite).First(&v).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				bs = append(bs, false)
			} else {
				return nil, err
			}
		} else {
			bs = append(bs, true)
		}
	}
	return bs, nil
}
