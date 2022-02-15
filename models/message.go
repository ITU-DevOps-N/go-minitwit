package models

import "time"

type Message struct {
	MessageID uint      `json:"messageID" gorm:"primaryKey"`
	AuthorID  uint      `json:"authorID"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Flagged   bool      `json:"flagged"`
}
