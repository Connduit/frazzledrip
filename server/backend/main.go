package main

import (
	"github.com/Connduit/frazzledrip/server/backend/handlers"
	"github.com/Connduit/frazzledrip/server/frontend"
	"log"
	"net/http"
)

func main() {
	// Serve static files
	http.Handle("/", frontend.Handler())

	// Register HTTP handlers
	handlers.RegisterHandlers()

	log.Println("C2 Server with Message Storage running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
