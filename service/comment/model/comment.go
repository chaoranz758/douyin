package model

import (
	"database/sql"
	"time"
)

type Commit struct {
	CommitID  int64        `gorm:"type:bigint(20);not null;uniqueIndex:commit;comment:'评论ID'" json:"commitID"`
	VideoID   int64        `gorm:"type:bigint(20);not null;index:video;comment:'该评论所属的视频ID'" json:"videoID"`
	UserID    int64        `gorm:"type:bigint(20);not null;index:user;comment:'评论创建用户ID'" json:"userID"`
	Message   string       `gorm:"type:varchar(8192);not null;comment:'评论内容'" json:"message"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt sql.NullTime `gorm:"index" json:"deletedAt"`
}
