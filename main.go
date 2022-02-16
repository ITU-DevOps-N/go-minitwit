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

	model "github.com/ITU-DevOps-N/go-minitwit/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func healtz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "Working good :)"})
}

func testDB() {

	// Create
	DB.Create(&model.User{Username: "emilravn", Email: "erav@itu.dk", Password: "123test"})
	DB.Create(&model.User{Username: "gianmarco", Email: "gimu@itu.dk", Password: "123test"})
	DB.Create(&model.User{Username: "tor", Email: "tor@itu.dk", Password: "123test"})
	DB.Create(&model.User{Username: "alex", Email: "alex@itu.dk", Password: "123test"})
	DB.Create(&model.User{Username: "henri", Email: "henri@itu.dk", Password: "123test"})

	DB.Create(&model.Message{Author: "emilravn", Text: "Hello World! My name is Emil Ravn"})
	DB.Create(&model.Message{Author: "gianmarco", Text: "Hello World! My name is Gianmarco"})
	DB.Create(&model.Message{Author: "tor", Text: "Hello World! My name is Tor"})
	DB.Create(&model.Message{Author: "alex", Text: "Hello World! My name is Alex"})
	DB.Create(&model.Message{Author: "henri", Text: "Hello World! My name is Henri"})

	// DB.Create(&model.Follow{Follower: 1, Following: 2})
	// DB.Create(&model.Follow{Follower: 1, Following: 3})
	// DB.Create(&model.Follow{Follower: 1, Following: 4})
	// DB.Create(&model.Follow{Follower: 1, Following: 5})

	// DB.Create(&model.Follow{Follower: 2, Following: 1})
	// DB.Create(&model.Follow{Follower: 2, Following: 3})
	// DB.Create(&model.Follow{Follower: 2, Following: 4})
	// DB.Create(&model.Follow{Follower: 2, Following: 5})

}

func SetupDB() {
	db, err := gorm.Open(sqlite.Open("minitwit.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("Failed to connect to database.")
	}

	// Migrate the entire table schema to the file according to our models
	db.AutoMigrate(&model.User{}, &model.Message{}, &model.Follow{})

	DB = db

	testDB()
}

func signUp (c *gin.Context) {
	var user model.User
	c.BindJSON(&user)
	DB.Create(&user)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func getUsers(c *gin.Context) {
	var users []model.User

	DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func getMessages(user string) []map[string]interface{} {
	var messages []model.Message

	DB.Find(&messages)
	var results []map[string]interface{}
	if user == "" {
		DB.Table("messages").Find(&results)
	} else {
		DB.Table("messages").Where("author = ?", user).Find(&results)
	}	
	return results
}

func timeline(c *gin.Context) {
	c.HTML(http.StatusOK, "timeline.tpl", gin.H{
		"title":    "Timeline",
		"endpoint": "public_timeline",
		"messages": getMessages(""),
	})
}

func user_timeline(c *gin.Context) {
	user := c.Request.URL.Query().Get("username")
	
	c.HTML(http.StatusOK, "timeline.tpl", gin.H{
		"title":    user + "'s Timeline",
		"endpoint": "timeline",
		"messages": getMessages(user),
	})
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tpl")
	r.Static("/static", "./static")

	SetupDB()
	r.GET("/info", healtz)
	r.GET("/", timeline)
	r.GET("/public_timeline", timeline)
	r.GET("/user_timeline", user_timeline)
	r.GET("/users", getUsers)
	r.GET("/messages", (func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getMessages("")})
	}))

	r.Run()
}
