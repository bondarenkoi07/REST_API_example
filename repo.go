package main

import (
	"app/REST_API_example/Controller"
	"log"
	"net/http"
	"time"
)

func main() {
	controller := Controller.NewController()

	srv := &http.Server{
		Handler: controller.Router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
