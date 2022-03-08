package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	// rename packe for easier read
	"net/http"

	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
)

func GetMessages(user string) []map[string]interface{} {
	var messages []model.Message

	database.DB.Find(&messages)
	var results []map[string]interface{}
	if user == "" {
		database.DB.Table("messages").Order("created_at desc").Find(&results)
	} else {
		database.DB.Table("messages").Order("created_at desc").Where("author = ?", user).Find(&results)
	}
	return results
}


func AddMessage(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	message := c.Request.FormValue("message")
	t := time.Now().Format(time.RFC822)
	time_now, _ := time.Parse(time.RFC822, t)
	database.DB.Create(&model.Message{Author: user, Text: message, CreatedAt: time_now})
	c.Redirect(http.StatusFound, "/user_timeline")
}



func GetFollower(follower uint, following uint) bool {
	var follows []model.Follow
	if follower == following {
		return false
	} else {
		database.DB.Find(&follows).Where("follower = ?", following).Where("following = ?", follower).First(&follows)
		return len(follows) > 0
	}
}
