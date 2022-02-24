package models

type LoginForm struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterForm struct {
	Username  string `form:"username" binding:"required"`
	Email     string `form:"email" binding:"required"`
	Password1 string `form:"password1" binding:"required"`
	Password2 string `form:"password2" binding:"required"`
}
