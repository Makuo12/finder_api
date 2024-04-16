package main

import (
	"finder_api/api"
	"log"
	"net/http"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Printf("Error creating server: %v\n", err)
	}
	server.Crawl("https://unsplash.com/")
	server.Receive()
	server.SetupRouter()
	server.Timer()
	err = http.ListenAndServe("localhost:8000", server.GetMux())
	if err != nil {
		log.Printf("Error setting up server: %v\n", err)
	} else {
		log.Println("Server running:")
	}
}
