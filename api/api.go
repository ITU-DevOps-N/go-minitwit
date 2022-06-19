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
	"os"
	"strconv"
	"strings"
	"time"

	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	bugsnaggin "github.com/bugsnag/bugsnag-go-gin"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
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
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_DATABASE") + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err1 := db.AutoMigrate(&model.User{}, &model.Message{}, &model.Follow{})
	if err1 != nil {
		panic(err1)
	}
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
	//If directories are unreferenced then they should be removed from the web root and/or the application directory.
	// reponse: HTTP/1.1 301 Moved Permanently

	Latest(c)
	var json model.RegisterForm

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// We have auto-complete here
	//Should maybe be fixed in the html: The first and most secure location is to disable the autocomplete attribute on the <form> HTML tag.

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

func GetMessages(user string, no int) []map[string]interface{} {
	var messages []model.Message

	DB.Find(&messages)
	var results []map[string]interface{}
	if user == "" {
		DB.Table("messages").Order("created_at desc").Limit(no).Find(&results)
	} else {
		DB.Table("messages").Order("created_at desc").Where("author = ?", user).Limit(no).Find(&results)
	}
	return results
}

func GetFollowers(user string) []string {
	var fllws []model.Follow
	var usr model.User
	followers := []string{}

	DB.Table("follows").Where("following = ?", GetUser(user).ID).Find(&fllws)
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

func sanitize(s string) string {
	return strings.ToValidUTF8(s, "")
}
func AddMessage(user string, message string) {
	//fix cross site forgery here
	//{
	// r := gin.Default()
	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession", store))
	// r.Use(csrf.Middleware(csrf.Options{
	// 	Secret: "secret123",
	// 	ErrorFunc: func(c *gin.Context) {
	// 		c.String(400, "CSRF token mismatch")
	// 		c.Abort()
	// 	},
	// }))
	message = sanitize(message)
	t := time.Now().Format(time.RFC822)
	time_now, _ := time.Parse(time.RFC822, t)
	DB.Create(&model.Message{Author: user, Text: message, CreatedAt: time_now})
}

func Latest(c *gin.Context) {
	l := c.Request.URL.Query().Get("latest")
	if l == "" {
		if c.FullPath() == "/latest" {
			c.JSON(200, gin.H{"latest": LATEST})
		}
		return
	}
	latest, err := strconv.Atoi(l)
	if err != nil {
		c.JSON(400, gin.H{"error_msg": "Latest must be an integer"})
		return
	}
	LATEST = latest
	if c.FullPath() == "/latest" {
		c.JSON(200, gin.H{"latest": LATEST})
	}

}

var cpuLoad = prometheus.NewGaugeFunc(prometheus.GaugeOpts{
	Name: "cpu_load_percentage",
	Help: "Current load of CPU in percentage",
}, getCpuLoad)

func getCpuLoad() float64 {
	cpuLoad, _ := cpu.Percent(time.Second, false)
	return cpuLoad[0]
}

func getGinMetrics(router *gin.Engine) {
	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/ginmetrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(router)
}

func init() {
	prometheus.MustRegister(cpuLoad)
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	router := gin.Default()
	router.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
		// Your Bugsnag project API key, required
		// set as environment variable $BUGSNAG_API_KEY
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "github.com/ITU-DevOps-N/go-minitwit"},
	}))
	SetupDB()
	//API ENDPOINTS ADDED
	router.GET("/", (func(c *gin.Context) {
		Latest(c)
		c.JSON(200, "Welcome to Go MiniTwit API!")
	}))

	router.GET("/version", (func(c *gin.Context) {
		Latest(c)
		c.Data(200, "application/json; charset=utf-8", []byte(os.Getenv("VERSION")))
	}))

	router.POST("/register", SignUp)

	// /msgs/*param means that param is optional
	// /msgs/:param means that param is required
	router.GET("/msgs/*usr", (func(c *gin.Context) {
		Latest(c)
		user := strings.Trim(c.Param("usr"), "/")
		no, err := strconv.Atoi(c.Request.URL.Query().Get("no"))
		if err != nil {
			no = 100
		}
		var data []map[string]interface{}
		if user == "" {
			data = GetMessages("", no)
		} else {
			data = GetMessages(user, no)
		}
		if len(data) == 0 {
			c.JSON(204, gin.H{})
		} else {
			c.JSON(http.StatusOK, gin.H{"data": data})
		}
	}))
	router.POST("/msgs/:usr", (func(c *gin.Context) {
		Latest(c)
		user := strings.Trim(c.Param("usr"), "/")

		if GetUser(user).ID == 0 {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		var message model.MessageForm

		if err := c.ShouldBindJSON(&message); err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error_msg": "You must provide a message"})
			return
		}

		AddMessage(user, message.Content)
		c.JSON(http.StatusNoContent, gin.H{})
	}))

	router.GET("/latest", Latest)

	router.GET("/fllws/:usr", (func(c *gin.Context) {
		Latest(c)
		user := strings.Trim(c.Param("usr"), "/")
		if GetUser(user).ID == 0 {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"data": GetFollowers(user)})
			return
		}
	}))

	router.POST("/fllws/:usr", (func(c *gin.Context) {
		Latest(c)
		user := strings.Trim(c.Param("usr"), "/")

		if GetUser(user).ID == 0 {
			c.JSON(404, gin.H{"error": "user not found"})
			return
		}

		var follow model.FollowForm
		if err := c.ShouldBindJSON(&follow); err != nil {
			c.JSON(403, gin.H{"error": err.Error()})
			return
		}
		if follow.Follow != "" {
			err := Follow(user, follow.Follow)
			if err != nil {
				c.JSON(403, gin.H{"error": ""})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
			return
		} else if follow.Unfollow != "" {
			err := Unfollow(user, follow.Unfollow)
			if err != nil {
				c.JSON(403, gin.H{"error": ""})
				return
			}
			c.JSON(http.StatusNoContent, gin.H{})
			return
		} else if len(follow.Latest) > 0 {
			latest, err := strconv.Atoi(follow.Latest[0])
			if err != nil {
				c.JSON(403, gin.H{"error_msg": "Latest must be an integer"})
				return
			}
			LATEST = latest
			c.JSON(http.StatusOK, gin.H{"data": GetFollowers(user)})
		} else {
			c.JSON(403, gin.H{"error_msg": "Only these fields are accepted: follow | unfollow | latest"})
			return
		}

	}))

	router.GET("/metrics", prometheusHandler())
	getGinMetrics(router)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
