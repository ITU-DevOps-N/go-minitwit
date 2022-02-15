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

	DB.Create(&model.Message{AuthorID: 1, Text: "Hello World! My name is Emil Ravn"})
	DB.Create(&model.Message{AuthorID: 2, Text: "Hello World! My name is Gianmarco"})
	DB.Create(&model.Message{AuthorID: 3, Text: "Hello World! My name is Tor"})
	DB.Create(&model.Message{AuthorID: 4, Text: "Hello World! My name is Alex"})
	DB.Create(&model.Message{AuthorID: 5, Text: "Hello World! My name is Henri"})

	user := model.User{}
	follower := model.User{}
	DB.Where("ID = ?", 2).First(&follower)
	DB.Where("ID = ?", 1).First(&user)
	user_followers := append(user.Followers, follower) // This still doesn't work

	// user.Followers = append(user.Followers, follower)
	// DB.Save(&user)
	DB.Model(&user).Update("Followers", user_followers)

}

func SetupDB() {
	db, err := gorm.Open(sqlite.Open("minitwit.db"), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("Failed to connect to database.")
	}

	// Migrate the entire table schema to the file according to our models
	db.AutoMigrate(&model.User{}, &model.Message{})

	DB = db

	testDB()

}
func getUsers(c *gin.Context) {
	var users []model.User

	DB.Find(&users)
	// model.db.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func getMessages(c *gin.Context) {
	var users []model.User

	DB.Find(&users)
	// model.db.Find(&books)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func timeline(c *gin.Context){
	DB.Where("message.flagged = 0 and message.author_id = user.user_id and (user.user_id = ? or user.user_id in (select whom_id from follower where who_id = ?))")

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

	SetupDB()
	r.GET("/info", healtz)
	r.GET("/", timeline)
	r.GET("/users", getUsers)
	r.GET("/messages", getMessages)
	
	r.Run()
}
