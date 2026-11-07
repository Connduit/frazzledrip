package main

import (
	"flag"
	"github.com/Connduit/frazzledrip/server/backend/handlers"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/frontend"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8081", "HTTP listen address")
	cacheDir := flag.String("cache-dir", network.DefaultDirectoryCacheDir(), "Directory snapshot cache folder")
	flag.Parse()
	if err := network.ConfigureDirectoryCache(*cacheDir); err != nil {
		log.Fatalf("Directory cache could not be loaded: %v", err)
	}

	// Serve static files
	http.Handle("/", frontend.Handler())

	// Register HTTP handlers
	handlers.RegisterHandlers()

	log.Printf("C2 Server with Message Storage starting on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
