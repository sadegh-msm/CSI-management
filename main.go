package main

import (
	"abrnoc_ch/handlers"
	"abrnoc_ch/routes"
	"database/sql"
	"fmt"
	"log"
)

type Server struct {
	Port string
	Host string
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "mydb"
)

// entry point of service
func main() {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		log.Fatal(err)
	}

	handlers.SetDB(db)

	defer db.Close()

	e := routes.Router()

	s := Server{
		Port: "8080",
		Host: "localhost",
	}

	log.Fatal(e.Start(s.Host + ":" + s.Port))
}
