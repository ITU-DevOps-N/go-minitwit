package main

/*
   Go MiniTwit
   ~~~~~~~~

   A microblogging application written in Golang with Gorilla.

   :copyright: (c) 2022 by DevØps - Group N.
   :license: BSD, see LICENSE for more details.
*/

import (
	"fmt"
	// "strings"
	// "html/template"
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
func getUsers(c *gin.Context) {
	var users []model.User

	DB.Find(&users)


	c.JSON(http.StatusOK, gin.H{"data": users})
}

func getMessages() []map[string]interface{} {
	var messages []model.Message

	DB.Find(&messages)
	var results []map[string]interface{}
	DB.Table("messages").Find(&results)

	return results
	// c.JSON(http.StatusOK, gin.H{"data": messages})
}

func timeline(c *gin.Context) {
	// c.JSON(http.StatusOK, gin.H{"data": getMessages()})
	c.HTML(http.StatusOK, "timeline.tpl", gin.H{
		"title":    "Timeline",
		"endpoint": "public_timeline",
		"messages": getMessages(),
	})

	//  query_db('''
	//     select message.*, user.* from message, user
	//     where message.flagged = 0 and message.author_id = user.user_id and (
	//         user.user_id = ? or
	//         user.user_id in (select whom_id from follower
	//                                 where who_id = ?))
	//     order by message.pub_date desc limit ?''',
	//     [session['user_id'], session['user_id'], PER_PAGE]))

}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tpl")
	r.Static("/css", "static")

	SetupDB()
	r.GET("/info", healtz)
	r.GET("/", timeline)
	r.GET("/public_timeline", timeline)
	r.GET("/users", getUsers)
	r.GET("/messages", (func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getMessages()})
		}))

	r.Run()
}
