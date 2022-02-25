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
	"net/http"
	"net/mail"
	"strings"
	"time"
	"strconv"

	model "github.com/ITU-DevOps-N/go-minitwit/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var LATEST = 0

func CreateUser(username string, email string, password string) bool {
	salt := Salt()
	err := DB.Create(&model.User{Username: username, Email: email, Salt: salt, Password: Hash(salt + password)}).Error
	if err != nil {
		return false
	} else {
		return true
	}
}

func Salt() string {
	bytes, _ := bcrypt.GenerateFromPassword(make([]byte, 8), 8)
	return string(bytes)
}

func Hash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func SetupDB() {
	db, err := gorm.Open(sqlite.Open("minitwit.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("Failed to connect to database.")
	}
	db.AutoMigrate(&model.User{}, &model.Message{}, &model.Follow{})
	DB = db
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidRegistration(c *gin.Context, username string, email string, password string) bool {
	if username == "" || email == "" || password == "" {
		c.JSON(400, gin.H{
			"error_msg": "All fields are required.",
		})
		return false
	}
	if !ValidEmail(email) {
		c.JSON(400, gin.H{
			"error_msg": "Email is not valid.",
		})
		return false
	}

	return true
}

func SignUp(c *gin.Context) {
	var json model.RegisterForm

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := json.Username
	email := json.Email
	password := json.Password

	if ValidRegistration(c, username, email, password) {
		fmt.Println(username)
		if CreateUser(username, email, password) {
			c.JSON(204, gin.H{})
		} else {
			c.JSON(400, gin.H{"error": "Username or email already exists."})
		}
	}
}

// PasswordCompare handles password hash compare
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

func GetUser(username string) model.User {
	var user model.User
	DB.Where("username = ?", username).First(&user) // SELECT * FROM USERS WHERE USERNAME = "?"
	return user
}

func GetMessages(user string) []map[string]interface{} {
	var messages []model.Message

	DB.Find(&messages)
	var results []map[string]interface{}
	if user == "" {
		DB.Table("messages").Order("created_at desc").Find(&results)
	} else {
		DB.Table("messages").Order("created_at desc").Where("author = ?", user).Find(&results)
	}
	return results
}

func GetFollowers(user string) []string {
	var fllws []model.Follow
	var usr model.User
	followers := []string{}
	
	DB.Table("follows").Where("following = ?", GetUser(user).ID ).Find(&fllws)
	for i := range fllws {
		DB.Where("ID = ?", fllws[i].Follower).First(&usr)
		followers = append(followers, usr.Username)
	}
	return followers
}

func GetFollower(follower uint, following uint) bool {
	var follows []model.Follow
	if follower == following {
		return false
	} else {
		DB.Find(&follows).Where("follower = ?", following).Where("following = ?", follower).First(&follows)
		return len(follows) > 0
	}
}

func Follow(user string, to_follow string) *gorm.DB {
	err := DB.Create(&model.Follow{Follower: GetUser(user).ID, Following: GetUser(to_follow).ID})	
	return err
}

func Unfollow(user string, to_unfollow string) *gorm.DB {
	var follows []model.Follow
	err := DB.Where("follower = ?", GetUser(user).ID).Where("following = ?", GetUser(to_unfollow).ID).Delete(&follows)
	return err
}

func AddMessage(user string, message string) {
	t := time.Now().Format(time.RFC822)
	time_now, _ := time.Parse(time.RFC822, t)
	DB.Create(&model.Message{Author: user, Text: message, CreatedAt: time_now})
}

func main() {
	router := gin.Default()

	SetupDB()
	//API ENDPOINTS ADDED
	router.GET("/", (func(c *gin.Context) {
		c.JSON(200, "Welcome to Go MiniTwit API!")
	}))
	router.POST("/register", SignUp)

	// /msgs/*param means that param is optional
	// /msgs/:param means that param is required
	router.GET("/msgs/*usr", (func(c *gin.Context) {
		user := strings.Trim(c.Param("usr"), "/")

		if user == "" {
			c.JSON(http.StatusOK, gin.H{"data": GetMessages("")})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": GetMessages(user)})
		}
	}))
	// messages_per_user (request method == POST) from minitwit_sim_api.py
	router.POST("/msgs/:usr", (func(c *gin.Context) {
		user := strings.Trim(c.Param("usr"), "/")

		var message model.MessageForm
		if user == "" {
			c.JSON(400, gin.H{"error_msg": "You must provide a username"})
		}

		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if message.Content == "" {
			c.JSON(400, gin.H{"error_msg": "You must provide a message"})
		}

		AddMessage(user, message.Content)

	}))

	// def update_latest(request: request):
	// 	global LATEST /latest?latest=
	// 	try_latest = request.args.get("latest", type=int, default=-1)
	// 	LATEST = try_latest if try_latest is not -1 else LATEST

	router.GET("/latest", func(c *gin.Context) {
		l := c.Request.URL.Query().Get("latest")
		if l == "" {
			l = "-1"
		}
		latest, err := strconv.Atoi(l)
		if err != nil {
			c.JSON(400, gin.H{"error_msg": "Latest must be an integer"})
			return
		}
		if latest  == -1{
			c.JSON(200, gin.H{"latest": LATEST})
		} else {
			LATEST = latest
			c.JSON(200, gin.H{"latest": LATEST})
		}
	})

	router.GET("/fllws/:usr", (func(c *gin.Context) {
		user := strings.Trim(c.Param("usr"), "/")
		if user == "" {
			c.JSON(400, gin.H{"error_msg": "You must provide a username"})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": GetFollowers(user)})
		}
	}))

	router.POST("/fllws/:usr", (func(c *gin.Context) {
		user := strings.Trim(c.Param("usr"), "/")
		if user == "" {
			c.JSON(400, gin.H{"error_msg": "You must provide a username"})
		}

		var follow model.FollowForm
		if err := c.ShouldBindJSON(&follow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if follow.Follow != "" {
			err := Follow(user, follow.Follow)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else if follow.Unfollow != "" {
			err := Unfollow(user, follow.Unfollow)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
		} else {
			c.JSON(400, gin.H{"error_msg": "You must provide a field to follow or unfollow"})
			return
		}

	}))
	router.Run(":8080")
}
