package handler

import (
	"encoding/json"
	"io"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/book"
	"log"
	"net/http"
	"strconv"
)

type DetailRequest struct {
	Isbn     string
	LibCodes []string
}

func GetBookRouter(query *sqlc.Queries) *http.ServeMux {
	bookRouter := http.NewServeMux()
	bookRouter.HandleFunc("POST /detail", func(w http.ResponseWriter, r *http.Request) {
		HandleBookRequests(w, r, query)
	})
	return bookRouter
}
func HandleBookRequests(w http.ResponseWriter, r *http.Request, query *sqlc.Queries) {
	var body DetailRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "Body decode error", http.StatusInternalServerError)
	}
	json.Unmarshal(b, &body)

	var LibCodeInt []int32

	for _, v := range body.LibCodes {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("err: %#+v\n", err)
			http.Error(w, "LibCode decode error", http.StatusInternalServerError)
		}
		LibCodeInt = append(LibCodeInt, int32(i))
	}

	data, err := book.RequestBookDetail(query, body.Isbn, LibCodeInt)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db requesting result error", http.StatusNotFound)
	}
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db data encoding error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
