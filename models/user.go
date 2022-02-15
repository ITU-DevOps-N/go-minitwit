package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" gorm:"type:varchar(100);not null"`
	Password  string    `json:"password" gorm:"not null"`
	// Messages  []Message `json:"messages" gorm:"foreignKey:AuthorID;"`       // A user can can have many messages
	Followers []User   `json:"followers" gorm:"many2many:user_followers;"` // Self-referencing many2many
}

type Message struct {
	MessageID uint      `json:"messageID" gorm:"primaryKey"`
	AuthorID  uint      `json:"authorID"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
	Flagged   bool      `json:"flagged"`
}