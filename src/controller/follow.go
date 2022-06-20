package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
)

func Follow(c *gin.Context) {
	user_to_follow := c.Request.URL.Query().Get("username")
	user, _ := c.Cookie("token")
	if user == "" {
		panic("You must be logged in to follow users")
	} else {
		database.DB.Create(&model.Follow{Follower: GetUser(user).ID, Following: GetUser(user_to_follow).ID})
	}
	c.Redirect(http.StatusFound, "/user_timeline?username="+user_to_follow)
}

func Unfollow(c *gin.Context) {
	var follows []model.Follow
	user_to_follow := c.Request.URL.Query().Get("username")
	user, _ := c.Cookie("token")
	if user == "" {
		panic("You must be logged in to follow users.")
	} else {
		database.DB.Where("follower = ?", GetUser(user).ID).Where("following = ?", GetUser(user_to_follow).ID).Delete(&follows)

	}
	c.Redirect(http.StatusFound, "/user_timeline?username="+user_to_follow)
}
