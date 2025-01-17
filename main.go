package main

import (
	"context"
	"flag"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/handler"
	"libra-backend/pkg/middleware"
	"log"
	"net/http"

	"github.com/dgraph-io/ristretto/v2"
)

var corsAllowList = []string{
	"http://127.0.0.1",
	"http://localhost:5173",
	"https://libsearch.xyz",
}

var cfg = config.GetEnvConfig()

func main() {
	port := flag.String("port", "3030", "default port 3030")
	deploy := flag.Bool("deploy", false, "default false")
	flag.Parse()

	log.Printf("Starting server on port: %s\n", *port)

	ctx := context.Background()

	var db_URL string
	if *deploy {
		db_URL = cfg.DATABASE_URL_SERVER
	} else {
		db_URL = cfg.DATABASE_URL
	}

	pool := db.ConnectPGPool(db_URL, ctx)
	defer pool.Close()

	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e3,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 28, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})

	if err != nil {
		panic(err)
	}

	router := http.NewServeMux()
	router.HandleFunc("/health", handler.GetHealth)

	scrapRouter := handler.GetScrapRouter(pool, cache)
	router.Handle("/scrap/", http.StripPrefix("/scrap", scrapRouter))

	bookRouter := handler.GetBookRouter(pool, cache)
	router.Handle("/book/", http.StripPrefix("/book", bookRouter))

	searchRouter, closeKiwi := handler.GetSearchRouter(pool)
	router.Handle("/search/", http.StripPrefix("/search", searchRouter))

	defer closeKiwi()
	if err := http.ListenAndServe(":"+*port, middleware.CorsMiddleware(router, corsAllowList)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
