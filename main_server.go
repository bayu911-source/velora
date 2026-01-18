package main

import (
	"log"

	"velora/server"
)

func main() {
	srv := server.NewServer()

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
