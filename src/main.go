package main

/*
   Go MiniTwit
   ~~~~~~~~
   A microblogging application written in Golang with Gin.
   :copyright: (c) 2022 by DevØps - Group N.
   :license: BSD, see LICENSE for more details.
*/

import (
	"fmt"
	"html/template"
	"time"

	// All of the below imports share the same package i.e. we could have
	// used the follow to access all functions.
	follow "github.com/ITU-DevOps-N/go-minitwit/src/controller"
	login "github.com/ITU-DevOps-N/go-minitwit/src/controller"

	messages "github.com/ITU-DevOps-N/go-minitwit/src/controller"
	registration "github.com/ITU-DevOps-N/go-minitwit/src/controller"
	ui "github.com/ITU-DevOps-N/go-minitwit/src/controller"
	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	"github.com/gin-gonic/gin"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", day, month, year, hour, minute, second)
}

func GetUserID(username string) uint {
	var user model.User
	database.DB.Where("username = ?", username).First(&user) // SELECT * FROM USERS WHERE USERNAME = "?"
	return user.ID
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
		"getUserId":    GetUserID,
	})
	router.LoadHTMLGlob("src/web/templates/*.tpl")
	router.Static("/web/static", "./src/web/static")

	database.SetupDB()
	router.GET("/", messages.Timeline)
	router.GET("/public_timeline", ui.Timeline)
	router.GET("/user_timeline", ui.UserTimeline)
	router.GET("/register", registration.Register)
	router.POST("/register", registration.SignUp)
	router.GET("/login", login.LoginPage)
	router.POST("/login", login.Login)
	router.GET("/logout", login.Logout)
	router.GET("/follow", follow.Follow)
	router.GET("/unfollow", follow.Unfollow)
	router.POST("/add_message", messages.AddMessage)

	router.Run(":81")
}
