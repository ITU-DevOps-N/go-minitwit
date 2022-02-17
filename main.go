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
	"net/url"

	"golang.org/x/crypto/bcrypt"

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
	createUser("emilravn", "erav@itu.dk", "123test")
	createUser("gianmarco", "gimu@itu.dk", "123test")
	createUser("tor", "tor@itu.dk", "123test")
	createUser("alex", "alex@itu.dk", "123test")
	createUser("henri", "henri@itu.dk", "123test")

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

func createUser(username string, email string, password string) {
	salt := salt()
	DB.Create(&model.User{Username: username, Email: email, Salt: salt, Password: hash(salt + password)})
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

func register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tpl", gin.H{
		"title":    "Sign Up",
		"endpoint": "user_signup",
	})
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func salt() string {
	bytes, _ := bcrypt.GenerateFromPassword(make([]byte, 8), 8)
	return string(bytes)
}
func hash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func signUp(c *gin.Context) {
	c.Request.ParseForm()
	username := c.Request.PostForm.Get("username")
	email := c.Request.PostForm.Get("email")
	password1 := c.Request.PostForm.Get("password")
	password2 := c.Request.PostForm.Get("password2")

	if password1 != password2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	if username == "" || email == "" || password1 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	if !validEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is not valid"})
		return
	}

	createUser(username, email, password1)
	// c.JSON(http.StatusOK, gin.H{"message": "User created"})
	location := url.URL{Path: "/login"}
	c.Redirect(http.StatusFound, location.RequestURI())

}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tpl", gin.H{
		"title":    "Login",
		"endpoint": "user_login",
	})
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
	// Should be changed to ReleaseMode when we are done coding
	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tpl")
	router.Static("/static", "./static")

	// Initialize
	SetupDB()

	// Endpoints
	router.GET("/", timeline)
	router.GET("/info", healtz)
	router.GET("/public_timeline", timeline)
	router.GET("/user_timeline", user_timeline)
	router.GET("/register", register)
	router.POST("/register", signUp)
	router.GET("/login", login)
	router.GET("/users", getUsers)
	router.GET("/messages", (func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getMessages("")})
	}))

	// Run the application
	router.Run()
}
