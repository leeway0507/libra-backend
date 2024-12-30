package main

import (
	"context"
	"flag"
	"libra-backend/config"
	"libra-backend/db"
	"libra-backend/db/sqlc"
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
	flag.Parse()
	log.Printf("Starting server on port: %s\n", *port)

	ctx := context.Background()
	conn := db.ConnectPG(cfg.DATABASE_URL, ctx)
	query := sqlc.New(conn)

	router := http.NewServeMux()
	router.HandleFunc("/health", handler.GetHealth)

	router.Handle("/static/", handler.StaticFileHandler())

	scrapRouter := handler.GetScrapRouter(query)
	router.Handle("/scrap/", http.StripPrefix("/scrap", scrapRouter))

	bookRouter := handler.GetBookRouter(query)
	router.Handle("/book/", http.StripPrefix("/book", bookRouter))

	searchRouter := handler.GetSearchRouter(query)
	router.Handle("/search/", http.StripPrefix("/search", searchRouter))

	if err := http.ListenAndServe(":"+*port, middleware.CorsMiddleware(router, corsAllowList)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
