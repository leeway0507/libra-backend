package handler

import (
	"log"
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Health check: %s %s", r.Method, r.URL.Path)
}
