package main

import (
	"flag"
	"github.com/Connduit/frazzledrip/server/backend/handlers"
	"github.com/Connduit/frazzledrip/server/frontend"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8081", "HTTP listen address")
	flag.Parse()

	// Serve static files
	http.Handle("/", frontend.Handler())

	// Register HTTP handlers
	handlers.RegisterHandlers()

	log.Printf("C2 Server with Message Storage starting on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
