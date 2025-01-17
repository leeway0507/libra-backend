package handler

import (
	"context"
	"encoding/json"
	"libra-backend/config"
	"libra-backend/db/sqlc"
	"libra-backend/pb"
	"libra-backend/pkg/bm25"
	"libra-backend/pkg/kiwi"
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

	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /normal", func(w http.ResponseWriter, r *http.Request) {
		HandleSearchQuery(w, r, pool, kb)
	})
	return searchRouter, kb.Close
}

func HandleSearchQuery(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, kb *kiwi.KiwiBuilder) {
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
	data, err := searchQuery.DBQuery().SearchFromBooks(ctx, sqlc.SearchFromBooksParams{
		Embedding: pgvector.NewVector(QueryResp.Embedding),
		LibCodes:  strings.Split(libCode, ","),
	})
	end := time.Now()
	log.Printf("query:%s time : %vms", keyword, end.Sub(start).Milliseconds())

	if err != nil {
		log.Printf("err: %#+v\n", err)
		http.Error(w, "db data encoding error", http.StatusInternalServerError)
		return
	}

	// RE_RANKING Based on Vector Search Result
	// Algorithm : OKAPI BM25
	if len(data) != 0 {
		// add keywords
		start := time.Now()
		keywords := strings.Split(keyword, " ")
		kb.AddWord(keyword, "NNP", 1)
		for i, k := range keywords {
			if k == keyword {
				break
			}
			score := 0.9 - float32(i)*0.2
			kb.AddWord(k, "NNP", max(0, score))
		}
		// load builder
		k := kb.Build()
		log.Println("kiwi load:", time.Since(start).String())
		defer k.Close()

		tokenizer := func(s string) []string {
			tokens, err := k.Analyze(s, 1, kiwi.KIWI_MATCH_ALL)
			if err != nil {
				log.Println("tokenizer error", err)
				return []string{}
			}
			return tokens
		}

		var corpus []string
		for _, d := range data {
			corpus = append(corpus, d.Title.String)
		}

		bm25, err := bm25.NewBM25Okapi(corpus, tokenizer, 2, 0.75, nil)
		if err != nil {
			log.Println("bm25:", err)
		}
		log.Println("keyword", keyword)
		tokenizedQuery := tokenizer(keyword)
		scores, err := bm25.GetScores(tokenizedQuery)
		if err != nil {
			log.Println(err)
		}
		// update score
		for i, v := range scores {
			data[i].Score = (1 - data[i].Score) + max(0, float32(v)*1.25)
		}
		// sort
		sort.Slice(data, func(i, j int) bool {
			return float64(data[i].Score) > float64(data[j].Score)
		})
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
