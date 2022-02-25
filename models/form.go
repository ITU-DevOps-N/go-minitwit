package models

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterForm struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"pwd" form:"password" binding:"required"`
}

type MessageForm struct {
	Content string `json:"content" form:"content" binding:"required"`
}

type FollowForm struct {
	Follow string `json:"follow" form:"follow"`
	Unfollow string `json:"unfollow" form:"follow"`
}
