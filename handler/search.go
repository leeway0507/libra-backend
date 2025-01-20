package handler

import (
	"context"
	"encoding/json"
	"libra-backend/config"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/kiwi"
	"libra-backend/pkg/pb"
	"libra-backend/pkg/search"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

var cfg = config.GetEnvConfig()

func GetSearchRouter(pool *pgxpool.Pool) (*http.ServeMux, func() int) {
	kiwi.Version()
	log.Println("tokenizer path :", cfg.TOKENIZER_PATH)
	kb := kiwi.NewBuilder(cfg.TOKENIZER_PATH, 1, kiwi.KIWI_BUILD_INTEGRATE_ALLOMORPH)
	k := kb.Build()

	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /normal", func(w http.ResponseWriter, r *http.Request) {
		HandleSearchQuery(w, r, pool, k)
	})
	return searchRouter, k.Close
}

func HandleSearchQuery(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, k *kiwi.Kiwi) {
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

	QueryResp, err := searchQuery.RequestQueryEmbedding(keyword)
	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "Failed to Request Keyword Embedding to OPEN AI", http.StatusInternalServerError)
		return
	}

	start := time.Now()
	keywords, err := k.Analyze_Noun(keyword, 1, kiwi.KIWI_MATCH_ALL)
	if err != nil {
		log.Printf("parseError: %#+v\n", err)
	}

	keywords = append(keywords, keyword)
	log.Printf("keywords: %#+v\n", keywords)
	var keywordForLike []string
	for _, v := range keywords {
		keywordForLike = append(keywordForLike, "%"+v+"%")
	}
	data, err := searchQuery.DBQuery().SearchFromBooks(ctx, sqlc.SearchFromBooksParams{
		Embedding: pgvector.NewVector(QueryResp.Embedding),
		Keywords:  keywordForLike,
		LibCodes:  strings.Split(libCode, ","),
	})
	end := time.Now()
	log.Printf("query:%s time : %vms", keyword, end.Sub(start).Milliseconds())

	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db data encoding error", http.StatusInternalServerError)
		return
	}
	// sort
	sort.Slice(data, func(i, j int) bool {
		return float64(data[i].Score) < float64(data[j].Score)
	})

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
