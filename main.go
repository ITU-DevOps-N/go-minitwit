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
	"strings"

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

	// testDB()
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

func validRegistration(c *gin.Context, username string, email string, password1 string, password2 string) bool {
	if password1 != password2 {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
            "error": "Passwords do not match",
		})
		return false
	}
	if username == "" || email == "" || password1 == "" {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
            "error": "All fields are required",
		})
		return false
	}
	if !validEmail(email) {
		c.HTML(http.StatusOK, "register.tpl", gin.H{
            "error": "Email is not valid",
		})
		return false
	}

	return true	
}

func signUp(c *gin.Context) {
	c.Request.ParseForm()
	username := c.Request.PostForm.Get("username")
	email := c.Request.PostForm.Get("email")
	password1 := c.Request.PostForm.Get("password")
	password2 := c.Request.PostForm.Get("password2")

	
	if validRegistration(c, username, email, password1, password2) {
		createUser(username, email, password1)
		// c.JSON(http.StatusOK, gin.H{"message": "User created"})
		location := url.URL{Path: "/login"}
		c.Redirect(http.StatusFound, location.RequestURI())
	}
}

// PasswordCompare handles password hash compare
func PasswordCompare(salt string, password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(salt+password))

	return err
}

func validUser(username string, password string) bool {
	
	user := getUser(username)
	
	if user.Username == "" {
		return false
	}
	
	err := PasswordCompare(user.Salt, password, user.Password)
	
	return err == nil
}

func login(c *gin.Context) {
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
	user := getUser(username)
	if validUser(username, password) {
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

func logout(c *gin.Context) {
    c.SetCookie("token", "", -1, "", "", false, true)

    c.Redirect(http.StatusTemporaryRedirect, "/")
}


func loginPage(c *gin.Context) {

	c.HTML(http.StatusOK, "login.tpl", gin.H{
		"title":    "username",
		"endpoint": "password",
	})
}

func getUsers(c *gin.Context) {
	var users []model.User

	DB.Find(&users)

	c.JSON(http.StatusOK, gin.H{"data": users})
}

func getUser(username string) model.User{
	var user model.User
	DB.Where("username = ?", username).First(&user)
	return user
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
	user, _ := c.Cookie("token")
	if user == "" {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":    "Timeline",
			"endpoint": "public_timeline",
			"messages": getMessages(""),
		})
	} else {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title": "Timeline",
			"user": user,
			"user_timeline": false,
			"messages": getMessages(""),
		})
	}
}

func getFollower(follower uint, following uint) bool{
	var follows []model.Follow
	DB.Find(&follows).Where("follower = ?", following).Where("following = ?", follower).First(&follows)
	return len(follows) > 0
}

func follow(c *gin.Context) {
	user_to_follow := c.Request.URL.Query().Get("username")
	user, _ := c.Cookie("token")
	if user == "" {
		panic("You must be logged in to follow users")
	} else {
		DB.Create(&model.Follow{Follower: getUser(user).ID, Following: getUser(user_to_follow).ID})
	}
	c.Redirect(http.StatusFound, "/user_timeline?username="+user_to_follow)
}

func unfollow(c *gin.Context) {
	var follows []model.Follow
	user_to_follow := c.Request.URL.Query().Get("username")
	user, _ := c.Cookie("token")
	if user == "" {
		panic("You must be logged in to follow users")
	} else {
		DB.Where("follower = ?", getUser(user).ID).Where("following = ?", getUser(user_to_follow).ID).Delete(&follows)

	}
	c.Redirect(http.StatusFound, "/user_timeline?username="+user_to_follow)
}

func user_timeline(c *gin.Context) {
	user_query := c.Request.URL.Query().Get("username")

	if user_query != "" {
		user, err := c.Cookie("token")
		if user != "" || err != nil {
			followed := getFollower(getUser(user_query).ID, getUser(user).ID)
			var user_page = false
			if user == user_query {
				user_page = true
			}
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title": user_query + "'s Timeline",
				"user_timeline": true,
				"private": true,
				"user": user_query,
				"followed": followed,
				"user_page": user_page,
				"messages": getMessages(user_query),
			})
		} else {

			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title": user_query + "'s Timeline",
				"user_timeline": true,
				"private": true,
				"messages": getMessages(user_query),
			})
		}
		
	} else {
		user, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/public_timeline")
		}
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title": "My Timeline",
			"user": user,
			"private": true,
			"user_page": true,
			"messages": getMessages(user),
		})
	}
}

func addMessage(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	message := c.Request.FormValue("message")
	DB.Create(&model.Message{Author: user, Text: message})


	c.Redirect(http.StatusFound, "/user_timeline")
}

func main() {
	// Should be changed to ReleaseMode when we are done coding
	// gin.SetMode(gin.DebugMode)

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
	router.GET("/login", loginPage)
	router.POST("/login", login)
	router.GET("/logout", logout)
	router.GET("/users", getUsers)
	router.GET("/follow", follow)
	router.GET("/unfollow", unfollow)
	router.POST("/add_message", addMessage)
	router.GET("/messages", (func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": getMessages("")})
	}))

	// Run the application
	router.Run()
}
