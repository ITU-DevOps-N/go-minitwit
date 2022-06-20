package controllers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(username string) model.User {
	var user model.User
	database.DB.Where("username = ?", username).First(&user) // SELECT * FROM USERS WHERE USERNAME = "?"
	return user
}

func PasswordCompare(salt string, password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(salt+password))
	return err
}

func ValidUser(username string, password string) bool {

	user := GetUser(username)

	if user.Username == "" {
		return false
	}

	err := PasswordCompare(user.Salt, password, user.Password)

	return err == nil
}

func Login(c *gin.Context) {
	c.Request.ParseForm()

	username := strings.ToLower(c.Request.PostForm.Get("username"))
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
			"ErrorTitle":   "Login failed.",
			"ErrorMessage": "Invalid credentials provided.",
		})
	}

	location := url.URL{Path: "/user_timeline"}
	c.Redirect(http.StatusFound, location.RequestURI())
}
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tpl", gin.H{
		"title": "Login",
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
