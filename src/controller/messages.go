package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	// rename packe for easier read
	"net/http"

	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
)

func GetMessages(user string, page string) []map[string]interface{} {
	var results []map[string]interface{}

	offset, messagesPerPage := LimitMessages(page)

	
	if user == "" {
		database.DB.Table("messages").Limit(messagesPerPage).Order("created_at desc").Offset(offset).Find(&results)
	} else {
		database.DB.Table("messages").Where("author = ?", user).Limit(messagesPerPage).Order("created_at desc").Offset(offset).Find(&results)
	}
	return results
}

func LimitMessages(page string) (int, int) {
	messagesPerPage := 50
	p, err := strconv.Atoi(page)
	if err != nil {
		panic("Failed to parse page number")
	}
	offset := (p - 1) * messagesPerPage
	return offset, messagesPerPage
}

func AddMessage(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	message := c.Request.FormValue("message")
	t := time.Now().Format(time.RFC1123)
	time_now, _ := time.Parse(time.RFC1123, t)
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
