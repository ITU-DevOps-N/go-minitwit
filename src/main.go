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
	"html/template"
	"time"

	// All of the below imports share the same package i.e. we could have
	// used the follow to access all functions.
	controller "github.com/ITU-DevOps-N/go-minitwit/src/controller"
	database "github.com/ITU-DevOps-N/go-minitwit/src/database"
	model "github.com/ITU-DevOps-N/go-minitwit/src/models"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
)

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

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", day, month, year, hour, minute, second)
}

func GetUserID(username string) uint {
	var user model.User
	database.DB.Where("username = ?", username).First(&user) // SELECT * FROM USERS WHERE USERNAME = "?"
	return user.ID
}

func main() {
	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
		"getUserId":    GetUserID,
	})

	router.LoadHTMLGlob("src/web/templates/*.tpl")
	router.Static("/web/static", "./src/web/static")

	database.SetupDB()
	router.GET("/", controller.Timeline)
	router.GET("/public_timeline", controller.Timeline)
	router.GET("/user_timeline", controller.UserTimeline)
	router.GET("/register", controller.Register)
	router.POST("/register", controller.SignUp)
	router.GET("/login", controller.LoginPage)
	router.POST("/login", controller.Login)
	router.GET("/logout", controller.Logout)
	router.GET("/follow", controller.Follow)
	router.GET("/unfollow", controller.Unfollow)
	router.POST("/add_message", controller.AddMessage)

	router.GET("/metrics", prometheusHandler())
	getGinMetrics(router)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
