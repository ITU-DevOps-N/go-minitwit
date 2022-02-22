package controller

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.Request.ParseForm()
	// session := sessions.Default(c)
	username := c.Request.PostForm.Get("username")
	password := c.Request.PostForm.Get("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusOK, "login.tpl", gin.H{
			"ErrorTitle":   "Empty Fields",
			"ErrorMessage": "Please fill in all fields",
		})
	}
	user := GetUser(username)
	if ValidUser(username, password) {
		// token := generateSessionToken()
		// session := sessions.Default(c)
		// session.Set("id", user.ID)
		// session.Set("email", user.Email)
		// session.Save()
		c.SetCookie("token", user.Username, 3600, "", "", false, true)

	} else {
		c.HTML(http.StatusOK, "login.tpl", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided",
		})
	}

	location := url.URL{Path: "/user_timeline"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)

	c.Redirect(http.StatusTemporaryRedirect, "/")
}
