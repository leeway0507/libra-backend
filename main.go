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
)

var corsAllowList = []string{
	"http://localhost:5173",
	"https://libra-client.pages.dev",
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

	router := http.NewServeMux()
	router.HandleFunc("/health", handler.GetHealth)

	router.Handle("/static/", handler.StaticFileHandler())

	scrapRouter := handler.GetScrapRouter(pool)
	router.Handle("/scrap/", http.StripPrefix("/scrap", scrapRouter))

	bookRouter := handler.GetBookRouter(pool)
	router.Handle("/book/", http.StripPrefix("/book", bookRouter))

	searchRouter := handler.GetSearchRouter(pool)
	router.Handle("/search/", http.StripPrefix("/search", searchRouter))

	if err := http.ListenAndServe(":"+*port, middleware.CorsMiddleware(router, corsAllowList)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
