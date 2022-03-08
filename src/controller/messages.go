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

func Timeline(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title": "Timeline",
			// "endpoint": "public_timeline",
			"messages": GetMessages(""),
		})
	} else {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":         "Timeline",
			"user":          user,
			"user_timeline": false,
			"messages":      GetMessages(""),
		})
	}
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

func UserTimeline(c *gin.Context) {
	user_query := c.Request.URL.Query().Get("username")

	if user_query != "" {
		user, err := c.Cookie("token")
		if user != "" || err != nil {
			followed := GetFollower(GetUser(user_query).ID, GetUser(user).ID)
			var user_page = false
			if user == user_query {
				user_page = true
			}
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title":         user_query + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"user":          user_query,
				"followed":      followed,
				"user_page":     user_page,
				"messages":      GetMessages(user_query),
			})
		} else {
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title":         user_query + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages(user_query),
			})
		}
	} else {
		user, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/public_timeline")
		}
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":     "My Timeline",
			"user":      user,
			"private":   true,
			"user_page": true,
			"messages":  GetMessages(user),
		})
	}
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