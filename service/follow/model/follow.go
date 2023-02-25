package model

type Follow struct {
	ID         int64 `gorm:"primarykey" json:"id"`
	FollowerID int64 `gorm:"type:bigint(20);not null;uniqueIndex:follow;comment:'关注者ID'" json:"followerID"`
	FollowID   int64 `gorm:"type:bigint(20);not null;uniqueIndex:follow;comment:'被关注者ID'" json:"followID"`
	Status     int8  `gorm:"not null;uniqueIndex:follow;default:0;comment:'关注或者删除'" json:"status"`
}
