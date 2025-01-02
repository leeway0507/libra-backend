package handler

import (
	"context"
	"encoding/json"
	"libra-backend/config"
	"libra-backend/db/sqlc"
	"libra-backend/pb"
	"libra-backend/pkg/search"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

var cfg = config.GetEnvConfig()

func GetSearchRouter(pool *pgxpool.Pool) *http.ServeMux {

	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /normal", func(w http.ResponseWriter, r *http.Request) {
		HandleSearchNormal(w, r, pool)
	})
	return searchRouter
}

func HandleSearchNormal(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := sqlc.New(conn)
	searchQuery := search.New(query, cfg.OPEN_AI_API_KEY)
	keyword := r.URL.Query().Get("q")
	if keyword == "" {
		log.Println("no query found", r.URL)
		http.Error(w, "no query found", http.StatusBadRequest)
		return
	}

	libCode := r.URL.Query().Get("libCode")
	if libCode == "" {
		log.Println("no query found", r.URL)
		http.Error(w, "no query found", http.StatusBadRequest)
		return
	}

	libCodeArr := strings.Split(libCode, ",")

	QueryResp, err := searchQuery.RequestQueryEmbedding(keyword)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "Failed to Request Keyword Embedding to OPEN AI", http.StatusInternalServerError)
		return
	}

	data, err := searchQuery.DBQuery().SearchFromBooks(ctx, sqlc.SearchFromBooksParams{
		Embedding: pgvector.NewVector(QueryResp.Embedding),
		LibCodes:  libCodeArr,
	})
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db data encoding error", http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db data encoding error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
	searchQuery.SaveQueryEmbedding(&pb.QueryEmbedding{
		Embedding: QueryResp.Embedding,
		Query:     QueryResp.Query,
	})
}
