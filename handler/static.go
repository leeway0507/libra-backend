package handler

import "net/http"

func StaticFileHandler() http.Handler {
	return http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
}
