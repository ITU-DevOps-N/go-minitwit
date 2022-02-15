package main

/*
   MiniTwit
   ~~~~~~~~

   A microblogging application written with Flask and sqlite3.

   :copyright: (c) 2010 by Armin Ronacher.
   :license: BSD, see LICENSE for more details.
*/

import (
	// "encoding/json"
	"fmt"
	"gorm.io/gorm"
	"net/http"

	"gorm.io/driver/sqlite"
	// _ "github.com/mattn/go-sqlite3"
)

// App struct to hold a database pointer
type App struct {
	DB *gorm.DB
}

// DB table structure
// TODO: field tags (primaryKeys etc.)
type follower struct {
	gorm.Model
	who_id  int
	whom_id int
}

type message struct {
	gorm.Model
	message_id int
	author_id  int
	text       string
	pub_date   int
	flagged    int
}

type user struct {
	gorm.Model
	user_id  int
	username string
	email    string
	pw_hash  string
}

type Timeline struct {
	gorm.Model
	Title    string
	Messages []message
}

func (a *App) Initialize(){
	db, err := gorm.Open(sqlite.Open("minitwit.db"), &gorm.Config{})

	// db, err := gorm.Open(sqlite.Open("minitwit.db"))
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		panic("failed to connect database")
	}

	a.DB = db

	// Migrate the schema.
	a.DB.AutoMigrate(&user{}, &follower{}, &message{}, &Timeline{})
}

func main() {
	a := &App{}
	a.Initialize()

	http.HandleFunc("/", a.handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

	defer a.DB.Close()
}

func (a *App) handler(w http.ResponseWriter, r *http.Request) {
	// // Create a test user.
	// a.DB.Create(&user{username: "Testname"})

	// // Read from DB.
	// var user1 user
	// a.DB.First(&user1, "username = ?", "Testname")

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte("hello world"))

	// // Delete.
	// a.DB.Delete(&user1)
}

/** GIAN MARCO **/
// // configuration

// var DATABASE = "./minitwit.db"
// var PER_PAGE = 30
// var DEBUG = true
// var SECRET_KEY = "development key"

// // create our little application
// var app = ""

// func main() {
// 	log.Println("starting API server")
// 	// create a new router
// 	router := mux.NewRouter()
// 	log.Println("creating routes")
// 	// specify endpoints
// 	router.HandleFunc("/", public_timeline)

// 	http.Handle("/", router)

// 	//start and listen to requests
// 	http.ListenAndServe(":8080", router)
// }

// func db() Timeline {
// 	db, err := sql.Open("sqlite3", DATABASE)
// 	checkErr(err)

// 	stmt, err := db.Prepare("select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?")
// 	checkErr(err)

// 	rows, err := stmt.Query(PER_PAGE)
// 	checkErr(err)

// 	var message_id int
// 	var author_id int
// 	var text string
// 	var pub_date int
// 	var flagged int

// 	var user_id int
// 	var username string
// 	var email string
// 	var pw_hash string

// 	var timeline Timeline
// 	timeline.Messages = make([]message, 0)
// 	for rows.Next() {
// 		err = rows.Scan(&message_id, &author_id, &text, &pub_date, &flagged, &user_id, &username, &email, &pw_hash) //(&uid, &username, &department, &created)
// 		checkErr(err)

// 		timeline.Messages = append(timeline.Messages, message{message_id, author_id, text, pub_date, flagged})
// 	}
// 	return timeline
// }

// func public_timeline(w http.ResponseWriter, r *http.Request) {
// 	// baseURL := *url.URL
// 	// fmap := template.FuncMap{
// 	// 	"url_for": func(path string) string {
// 	// 		return baseURL.String() + path
// 	// 	},

// 	// }
// 	template.Must(template.ParseFiles("templates/timeline-test.html")).Execute(w, db())
// }

// func checkErr(err error) {
// 	if err != nil {
// 		panic("failed to connect database " + err.Error())
// 	}
// }
