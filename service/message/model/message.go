package model

import (
	"database/sql"
	"time"
)

type Messages []Message

type Message struct {
	MessageID int64        `gorm:"type:bigint(20);not null;unique;comment:'消息ID'" json:"messageID"`
	UserID    int64        `gorm:"type:bigint(20);not null;index:user;comment:'发消息的用户ID'" json:"userID"`
	ToUserID  int64        `gorm:"type:bigint(20);not null;index:user;comment:'接收消息的用户ID'" json:"toUserID"`
	Content   string       `gorm:"type:varchar(8192);not null;comment:'消息内容'" json:"content"`
	CreatedAt int64        `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
	DeletedAt sql.NullTime `gorm:"index" json:"deletedAt"`
}

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].CreatedAt < m[j].CreatedAt
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i] //与上面的交换方法一致
}
