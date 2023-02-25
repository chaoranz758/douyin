package model

import (
	"database/sql"
	"time"
)

type Video struct {
	VideoID       int64        `gorm:"primarykey;type:bigint(20);not null;uniqueIndex:video;comment:'视频ID'" json:"videoID"`
	UserID        int64        `gorm:"type:bigint(20);not null;index:user;comment:'发布该视频的用户ID'" json:"userID"`
	Commit        int64        `gorm:"type:bigint(10);not null;default:0;comment:'视频评论数'" json:"commit"`
	FavoriteVideo int64        `gorm:"type:bigint(10);not null;default:0;comment:'视频点赞数'" json:"favoriteVideo"`
	VideoUrl      string       `gorm:"type:varchar(500);not null;comment:'视频地址'" json:"videoUrl"`
	CoverUrl      string       `gorm:"type:varchar(500);not null;comment:'视频封面地址'" json:"coverUrl"`
	Title         string       `gorm:"type:varchar(200);not null;comment:'视频标题'" json:"title"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	DeletedAt     sql.NullTime `gorm:"index" json:"deletedAt"`
}
