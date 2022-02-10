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
	// "fmt"
	"log"
    "net/url"
	"html/template"
	"net/http"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
)

// DB table structure
type follower struct {
	who_id  int
	whom_id int
}

type message struct {
	message_id int
	author_id  int
	text       string
	pub_date   int
	flagged    int
}

type user struct {
	user_id  int
	username string
	email    string
	pw_hash  string
}

type Timeline struct {
    Title       string
    Messages    []message
}

// configuration

var DATABASE = "./minitwit.db"
var PER_PAGE = 30
var DEBUG = true
var SECRET_KEY = "development key"

// create our little application
var app = ""

func main() {
	log.Println("starting API server")
	// create a new router
	router := mux.NewRouter()
	log.Println("creating routes")
	// specify endpoints
	router.HandleFunc("/", public_timeline)

	http.Handle("/", router)

	//start and listen to requests
	http.ListenAndServe(":8080", router)
}

func db() Timeline {
	db, err := sql.Open("sqlite3", DATABASE)
	checkErr(err)

    stmt, err := db.Prepare("select message.*, user.* from message, user where message.flagged = 0 and message.author_id = user.user_id order by message.pub_date desc limit ?")
    checkErr(err)

    rows, err := stmt.Query(PER_PAGE)
    checkErr(err)

	var message_id int
	var author_id  int
	var text       string
	var pub_date   int
	var flagged    int

	var user_id  int
	var username string
	var email    string
	var pw_hash  string

    var timeline Timeline
    timeline.Messages = make([]message, 0)
	for rows.Next() {
		err = rows.Scan(&message_id, &author_id, &text, &pub_date, &flagged, &user_id, &username, &email, &pw_hash) //(&uid, &username, &department, &created)
		checkErr(err)
		
        timeline.Messages = append(timeline.Messages, message{message_id, author_id, text, pub_date, flagged})
    }
    return timeline
}

func public_timeline(w http.ResponseWriter, r *http.Request) {
    // baseURL := *url.URL
    // fmap := template.FuncMap{
	// 	"url_for": func(path string) string {
	// 		return baseURL.String() + path
	// 	},

	// }
    template.Must(template.ParseFiles("templates/timeline-test.html")).Execute(w, db())
}

func checkErr(err error) {
	if err != nil {
		panic("failed to connect database " + err.Error())
	}
}
