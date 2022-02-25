package api

/*
   Go MiniTwit
   ~~~~~~~~
   A microblogging application written in Golang with Gin.
   :copyright: (c) 2022 by DevØps - Group N.
   :license: BSD, see LICENSE for more details.
*/

import (
	"fmt"
	"html/template"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
	"time"

	model "github.com/ITU-DevOps-N/go-minitwit/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func CreateUser(username string, email string, password string) {
	salt := Salt()
	DB.Create(&model.User{Username: username, Email: email, Salt: salt, Password: Hash(salt + password)})
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

func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tpl", gin.H{
		"title": "Register",
	})
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
			"ErrorTitle":   "Login failed.",
			"ErrorMessage": "Invalid credentials provided.",
		})
	}

	location := url.URL{Path: "/user_timeline"}
	c.Redirect(http.StatusFound, location.RequestURI())
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tpl", gin.H{
		"title": "Login",
	})
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

func Timeline(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title": "Timeline",
			// "endpoint": "public_timeline",
			"messages": GetMessages(""),
		})
	} else {
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":         "Timeline",
			"user":          user,
			"user_timeline": false,
			"messages":      GetMessages(""),
		})
	}
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

func Follow(c *gin.Context) {
	user_to_follow := c.Request.URL.Query().Get("username")
	user, _ := c.Cookie("token")
	if user == "" {
		panic("You must be logged in to follow users")
	} else {
		DB.Create(&model.Follow{Follower: GetUser(user).ID, Following: GetUser(user_to_follow).ID})
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
		DB.Where("follower = ?", GetUser(user).ID).Where("following = ?", GetUser(user_to_follow).ID).Delete(&follows)

	}
	c.Redirect(http.StatusFound, "/user_timeline?username="+user_to_follow)
}

func UserTimeline(c *gin.Context) {
	user_query := c.Request.URL.Query().Get("username")

	if user_query != "" {
		user, err := c.Cookie("token")
		if user != "" || err != nil {
			followed := GetFollower(GetUser(user_query).ID, GetUser(user).ID)
			var user_page = false
			if user == user_query {
				user_page = true
			}
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title":         user_query + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"user":          user_query,
				"followed":      followed,
				"user_page":     user_page,
				"messages":      GetMessages(user_query),
			})
		} else {
			c.HTML(http.StatusOK, "timeline.tpl", gin.H{
				"title":         user_query + "'s Timeline",
				"user_timeline": true,
				"private":       true,
				"messages":      GetMessages(user_query),
			})
		}
	} else {
		user, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, "/public_timeline")
		}
		c.HTML(http.StatusOK, "timeline.tpl", gin.H{
			"title":     "My Timeline",
			"user":      user,
			"private":   true,
			"user_page": true,
			"messages":  GetMessages(user),
		})
	}
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d", day, month, year, t.Hour(), t.Minute())
}

func GetUserID(username string) uint {
	var user model.User
	DB.Where("username = ?", username).First(&user) // SELECT * FROM USERS WHERE USERNAME = "?"
	return user.ID
}

func AddMessage(c *gin.Context) {
	user, _ := c.Cookie("token")
	if user == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
	message := c.Request.FormValue("message")
	t := time.Now().Format(time.RFC822)
	time_now, _ := time.Parse(time.RFC822, t)
	DB.Create(&model.Message{Author: user, Text: message, CreatedAt: time_now})
	c.Redirect(http.StatusFound, "/user_timeline")
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
		"getUserId":    GetUserID,
	})
	router.LoadHTMLGlob("templates/*.tpl")
	router.Static("/static", "./static")

	SetupDB()
	router.GET("/", Timeline)
	router.GET("/public_timeline", Timeline)
	router.GET("/user_timeline", UserTimeline)
	router.GET("/register", Register)
	router.POST("/register", SignUp)
	router.GET("/login", LoginPage)
	router.POST("/login", Login)
	router.GET("/logout", Logout)
	// router.GET("/users", GetUsers)
	router.GET("/follow", Follow)
	router.GET("/unfollow", Unfollow)
	router.POST("/add_message", AddMessage)
	//API ENDPOINTS ADDED

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
	router.POST("/msgs/:usr", (func(c *gin.Context) {
		user := strings.Trim(c.Param("usr"), "/")
		// TODO: continue from the 'messages_per_user(username)' method in API (line 146)

	}))

	router.GET("/latest", func(c *gin.Context) {
		_, val := c.Get("latest")
		c.JSON(http.StatusOK, gin.H{"data": val})

	})

	router.GET("/msgs", func(c *gin.Context) {

	})

	router.Run()
}

// Helper method for API
func notReqFromSimulator()
