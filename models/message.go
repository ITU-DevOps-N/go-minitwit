package model

import "time"

type Message struct {
	MessageID uint `gorm:"primaryKey"`
	UserID    uint // A message belongs to a single user
	Text      string
	CreatedAt time.Time
	Flagged   int
}