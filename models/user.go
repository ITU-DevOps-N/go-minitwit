package models

type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique;type:varchar(100);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);not null"`
	Salt     string `json:"salt" gorm:"type:varchar(100);not null"`
	Password string `json:"password" gorm:"not null"`
}

type Message struct {
	MessageID uint   `json:"messageID" gorm:"primaryKey"`
	Author    string `json:"authorID" gorm:"foreignKey:Username;"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	Flagged   bool   `json:"flagged"`
}

type Follow struct {
	Follower  uint `json:"follower" gorm:"primaryKey;foreignKey:ID;"`
	Following uint `json:"following" gorm:"primaryKey;foreignKey:ID;"`
}
