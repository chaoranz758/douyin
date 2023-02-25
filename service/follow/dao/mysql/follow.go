package mysql

import (
	"douyin/proto/follow/request"
	"douyin/service/follow/model"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	follow = iota
	cancelFollow
)

const (
	errorFollowRelation0 = "关注关系数据没有插入进follow表但没有报错"
	errorFollowRelation  = "create follow relation failed"
	errorDeleteFollow0   = "关注关系在follow表中没有被删除成功"
	errorDeleteFollow    = "delete follow failed"
)

var (
	errorFollowRelation0V = errors.New("关注关系数据没有插入进follow表但没有报错")
	errorDeleteFollow0V   = errors.New("关注关系在follow表中没有被删除成功")
)

func CreateFollow(f *model.Follow) error {
	var fo model.Follow
	resultGet := db.Where("follower_id = ? and follow_id = ? and status = ?", f.FollowerID, f.FollowID, cancelFollow).First(&fo)
	if resultGet.RowsAffected == 0 {
		result := db.Create(f)
		if result.RowsAffected == 0 {
			zap.L().Error(errorFollowRelation0, zap.Error(errorFollowRelation0V))
			return errorFollowRelation0V
		}
		if result.Error != nil {
			zap.L().Error(errorFollowRelation, zap.Error(result.Error))
			return result.Error
		}
		return nil
	}
	if resultGet.Error != nil {
		zap.L().Error(errorFollowRelation, zap.Error(resultGet.Error))
		return resultGet.Error
	}
	zap.L().Info("用户取消关注后重新关注")
	result := db.Model(&model.Follow{}).Where("follower_id = ? and follow_id = ? and status = ?", f.FollowerID, f.FollowID, cancelFollow).Update("status", follow)
	if result.RowsAffected == 0 {
		zap.L().Error(errorFollowRelation0, zap.Error(errorFollowRelation0V))
		return errorFollowRelation0V
	}
	if result.Error != nil {
		zap.L().Error(errorFollowRelation, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func DeleteFollow(userId, toUserId int64) error {
	result := db.Model(&model.Follow{}).Where("follower_id = ? and follow_id = ? and status = ?", userId, toUserId, follow).Update("status", cancelFollow)
	if result.RowsAffected == 0 {
		zap.L().Error(errorDeleteFollow0, zap.Error(errorDeleteFollow0V))
		return errorDeleteFollow0V
	}
	if result.Error != nil {
		zap.L().Error(errorDeleteFollow, zap.Error(result.Error))
		return result.Error
	}
	return nil
}

func GetUserFollowList(fs *[]model.Follow, userId int64) error {
	if err := db.Where("follower_id = ? and status = ?", userId, follow).Find(fs).Error; err != nil {
		return err
	}
	return nil
}

func GetUserFollowerList(fs *[]model.Follow, userId int64) error {
	if err := db.Where("follow_id = ? and status = ?", userId, follow).Find(fs).Error; err != nil {
		return err
	}
	return nil
}

func JudgeUserIsFollow(req *request.DouyinGetFollowRequest) (bool, error) {
	var f []model.Follow
	result := db.Where("follower_id = ? and follow_id = ? and status = ?", req.UserId, req.ToUserId, follow).Find(&f)
	if result.RowsAffected == 0 {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func JudgeUserIsFollowList(req *request.DouyinGetFollowListRequest) ([]bool, error) {
	var bs []bool
	for i := 0; i < len(req.ToUserId); i++ {
		var f model.Follow
		if err := db.Where("follower_id = ? and follow_id = ?  and status = ?", req.UserId, req.ToUserId[i], follow).First(&f).Error; err != nil {
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
