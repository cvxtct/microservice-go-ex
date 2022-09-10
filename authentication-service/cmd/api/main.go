package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// TODO connect to DB

	// set up config

	app := Config{}

	// define server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// start server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}
