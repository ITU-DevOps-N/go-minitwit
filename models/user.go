package model

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  string
	Messages  []Message  // A user can can have many messages
	Followers []Follower `gorm:"many2many:user_followers;"`
}