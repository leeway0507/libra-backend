package main

import (
	"flag"
	"libra-backend/handler"
	"libra-backend/pkg/middleware"
	"log"
	"net/http"
)

var corsAllowList = []string{
	"http://localhost:5173",
	"https://libra-client.pages.dev",
}

func main() {
	port := flag.String("port", "3030", "default port 3030")
	flag.Parse()
	log.Printf("Starting server on port: %s\n", *port)

	router := http.NewServeMux()
	router.HandleFunc("/health", handler.GetHealth)
	router.HandleFunc("/scrap/{libCode}/{isbn}", handler.HandleScrap)
	router.Handle("/static/", handler.StaticFileHandler())

	if err := http.ListenAndServe(":"+*port, middleware.CorsMiddleware(router, corsAllowList)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
