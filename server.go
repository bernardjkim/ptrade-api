package main

import (
	"flag"
	"log"
	"os"

	"github.com/bkim0128/bjstock-rest-service/src/system/app"
	"github.com/joho/godotenv"

	DB "github.com/bkim0128/bjstock-rest-service/src/system/db"
)

var (
	port      string
	dbURL     string
	dboptions string
)

func init() {
	flag.StringVar(&port, "port", "8080", "Accepting the port that the server should listen on")
	flag.StringVar(&dboptions, "dboptions", "parseTime=true", "Set the port for the application")

	flag.Parse()

	_ = godotenv.Load()

	if url := os.Getenv("JAWSDB_URL"); len(url) > 0 {
		dbURL = url
	}

	envPort := os.Getenv("PORT")
	if len(envPort) > 0 {
		port = envPort
	}
}

func main() {

	db, err := DB.ConnectURL(dbURL, dboptions)
	if err != nil {
		log.Println("Unable to connect to db")
		panic(err)
	}

	// DB.Init(db)

	s := app.NewServer()
	s.Init(port, db)
	s.Start()

}
