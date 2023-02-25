package model

import (
	"database/sql"
	"time"
)

type User struct {
	UserID          int64        `gorm:"primarykey;type:bigint(20);not null;uniqueIndex:userid;comment:'用户ID'" json:"userID"`
	UserName        string       `gorm:"type:varchar(32);not null;uniqueIndex:usernames;comment:'用户名'" json:"userName"`
	Password        string       `gorm:"type:varchar(32);not null;comment:'用户密码'" json:"password"`
	Follow          int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户关注数'" json:"follow"`
	Follower        int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户粉丝数'" json:"follower"`
	Friend          int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户朋友数'" json:"friend"`
	Video           int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户发布视频数'" json:"video"`
	FavoriteVideo   int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户点赞视频数'" json:"favoriteVideo"`
	FavoriteTotal   int64        `gorm:"type:bigint(10);not null;default:0;comment:'用户被赞总数'" json:"favoriteTotal"`
	Avatar          string       `gorm:"type:varchar(200);comment:'用户头像'" json:"avatar"`
	BackgroundImage string       `gorm:"type:varchar(200);comment:'用户个人页顶部大图'" json:"backgroundImage"`
	Signature       string       `gorm:"type:varchar(200);comment:'个人简介'" json:"signature"`
	CreatedAt       time.Time    `json:"createdAt"`
	UpdatedAt       time.Time    `json:"updatedAt"`
	DeletedAt       sql.NullTime `gorm:"index" json:"deletedAt"`
}
