package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Text string `json:"text,omitempty"`
	ToUserID uint `json:"toUserID,omitempty"`
	FromUserID uint `json:"fromUserID,omitempty"`

	ToUser User `json:"toUser,omitempty" gorm:"foreignKey:ToUserID"`
	FromUser User `json:"romUser,omitempty" gorm:"foreignKey:FromUserID"`
}