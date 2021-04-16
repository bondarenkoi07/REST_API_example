package main

import (
	"app/REST_API_example/Controller"
	"log"
	"net/http"
	"time"
)

func main() {
	controller := Controller.NewController()

	log.Print("Hello, World!")

	srv := &http.Server{
		Handler: controller.Router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
