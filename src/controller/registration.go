package controllers

import (
	"github.com/gin-gonic/gin"
	// rename packe for easier read
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username string, email string, password string) {
	salt := Salt()
	usr := strings.ToLower(username)
	database.DB.Create(&model.User{Username: usr, Email: email, Salt: salt, Password: Hash(salt + password)})
}

func Salt() string {
	bytes, _ := bcrypt.GenerateFromPassword(make([]byte, 8), 8)
	return string(bytes)
}

func Hash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidRegistration(c *gin.Context, username string, email string, password1 string, password2 string) bool {
	if password1 != password2 {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
			"error": "Passwords do not match.",
		})
		return false
	}
	if username == "" || email == "" || password1 == "" {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
			"error": "All fields are required.",
		})
		return false
	}
	if !ValidEmail(email) {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
			"error": "Email is not valid.",
		})
		return false
	}

	return true
}

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tpl", gin.H{
		"title": "Register",
	})
}

func SignUp(c *gin.Context) {
	c.Request.ParseForm()
	username := c.Request.PostForm.Get("username")
	email := c.Request.PostForm.Get("email")
	password1 := c.Request.PostForm.Get("password1")
	password2 := c.Request.PostForm.Get("password2")

	if ValidRegistration(c, username, email, password1, password2) {
		CreateUser(username, email, password1)
		location := url.URL{Path: "/login"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}
