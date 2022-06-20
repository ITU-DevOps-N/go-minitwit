package controllers

import (
	"github.com/gin-gonic/gin"
	// rename packe for easier read
	"net/http"
)

func Timeline(c *gin.Context) {
	user, _ := c.Cookie("token")

	page := c.DefaultQuery("page", "0")

	if user == "" {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title": "Timeline",
			// "endpoint": "public_timeline",
			"messages": GetMessages("", page),
		})
	} else {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":         "Timeline",
			"user":          user,
			"user_timeline": false,
			"messages":      GetMessages("", page),
		})
	}
}

func UserTimeline(c *gin.Context) {
	user_query := c.Request.URL.Query().Get("username")

	page := c.DefaultQuery("page", "0")

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
				"messages":      GetMessages(user_query, page),
			})
		} else {
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title":         user_query + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages(user_query, page),
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
			"messages":  GetMessages(user, page),
		})
	}
}
