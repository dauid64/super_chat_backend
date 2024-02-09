package models

type MessageWsData struct {
	Text string `json:"text,omitempty"`
	ToUserID uint `json:"toUserID,omitempty"`
	FromUserID uint `json:"fromUserID,omitempty"`
}