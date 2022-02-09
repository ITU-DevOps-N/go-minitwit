package main

/*
    MiniTwit
    ~~~~~~~~

    A microblogging application written with Flask and sqlite3.

    :copyright: (c) 2010 by Armin Ronacher.
    :license: BSD, see LICENSE for more details.
*/

import ("fmt"
		"sqlite3"
)



// configuration
var DATABASE = "/tmp/minitwit.db"
var PER_PAGE = 30
var DEBUG = true
var SECRET_KEY = "development key"

// create our little application
var app = ""


func main() {
	fmt.Println("Hello, World")
}
