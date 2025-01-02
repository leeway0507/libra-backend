package handler

import (
	"context"
	"encoding/json"
	"io"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/book"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DetailRequest struct {
	Isbn     string
	LibCodes []string
}

func GetBookRouter(pool *pgxpool.Pool) *http.ServeMux {
	bookRouter := http.NewServeMux()
	bookRouter.HandleFunc("POST /detail", func(w http.ResponseWriter, r *http.Request) {
		HandleBookRequests(w, r, pool)
	})
	return bookRouter
}
func HandleBookRequests(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := sqlc.New(conn)
	var body DetailRequest
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "Body decode error", http.StatusInternalServerError)
	}
	json.Unmarshal(b, &body)

	data, err := book.RequestBookDetail(query, body.Isbn, body.LibCodes)
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
