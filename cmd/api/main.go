package main

import (
	"log"
	"net/http"
	"notifications/controllers"
)

func main() {
	log.Printf("Listening to http://0.0.0.0:8080/")

	server := http.Server{
		Addr:    ":8080",
		Handler: controllers.NewHTTPHandler(),
	}

	log.Fatal(server.ListenAndServe())
}
