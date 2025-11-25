package main

import (
	"backend/handlers"
	"log"
	"net/http"
)

func main() {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir(".")))

	// Register HTTP handlers
	handlers.RegisterHandlers()

	log.Println("C2 Server with Message Storage running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
