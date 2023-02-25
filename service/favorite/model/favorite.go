package model

type Favorite struct {
	ID      int64 `gorm:"primarykey" json:"id"`
	VideoID int64 `gorm:"type:bigint(20);not null;uniqueIndex:favorite;comment:'用户点赞的视频ID'" json:"videoID"`
	UserID  int64 `gorm:"type:bigint(20);not null;uniqueIndex:favorite;comment:'用户ID'" json:"userID"`
	Status  int8  `gorm:"not null;uniqueIndex:favorite;default:0;comment:'点赞或者删除'" json:"status"`
}
