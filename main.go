package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request received: %s %s", r.Method, r.URL.Path)
		fmt.Fprintln(w, "hello world!!!")
	})

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Printf("err: %#+v\n", err)
	}
}
