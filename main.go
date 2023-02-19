package main

import (
	"abrnoc_ch/handlers"
	"abrnoc_ch/routes"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type Server struct {
	Port string
	Host string
}

var db *sql.DB

// entry point of service
func main() {
	var err error
	db, err = sql.Open("postgres", "postgres://user:password@localhost/dbname?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	handlers.SetDB(db)

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	e := routes.Router()

	e.Logger.Fatal(e.Start(":8080"))
}
