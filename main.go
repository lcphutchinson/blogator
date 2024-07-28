package main

import _ "github.com/lib/pq"
import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lcphutchinson/database"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("CONN")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	config := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	bootMux(config, mux)

	rssWorker := worker{
		DB: 		dbQueries,
		batchSize:	5,
		loopInterval:	30 * time.Second,
	}
	go rssWorker.Work()

	server := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}
	log.Fatal(server.ListenAndServe())
}
