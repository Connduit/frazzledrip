package main

import (
	"flag"
	"github.com/Connduit/frazzledrip/server/backend/handlers"
	"github.com/Connduit/frazzledrip/server/backend/network"
	"github.com/Connduit/frazzledrip/server/backend/state"
	"github.com/Connduit/frazzledrip/server/frontend"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8081", "HTTP listen address")
	cacheDir := flag.String("cache-dir", network.DefaultDirectoryCacheDir(), "Directory snapshot cache folder")
	database := flag.String("database", "", "Persistent server database (default: cache-dir/server.db)")
	flag.Parse()
	if *database == "" {
		*database = filepath.Join(*cacheDir, "server.db")
	}
	store, err := state.Open(*database)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	if err := network.ConfigurePersistentState(store); err != nil {
		log.Fatal(err)
	}
	if err := network.ConfigureDirectoryCache(*cacheDir); err != nil {
		log.Fatalf("Directory cache could not be loaded: %v", err)
	}

	// Serve static files
	http.Handle("/", frontend.Handler())

	// Register HTTP handlers
	handlers.RegisterHandlers()

	log.Printf("C2 Server with Message Storage starting on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, handlers.AuditHTTP(http.DefaultServeMux)))
}
