package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/book"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DetailRequest struct {
	Isbn     string
	LibCodes []string
}

func GetBookRouter(pool *pgxpool.Pool, cache *ristretto.Cache[string, []byte]) *http.ServeMux {
	bookRouter := http.NewServeMux()
	bookRouter.HandleFunc("POST /detail", func(w http.ResponseWriter, r *http.Request) {
		HandleBookRequests(w, r, pool)
	})

	// best sellers
	bestSeller := book.NewBestSeller(cfg.ALADIN_API_KEY)
	for _, cat := range bestSeller.GetCatName() {
		bookRouter.HandleFunc(
			fmt.Sprintf("GET /bestseller/%s", cat),
			func(w http.ResponseWriter, r *http.Request) {
				HandleBestSellerRequests(w, r, cat, bestSeller.Instance(cat), pool, cache)
			})
	}
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

func HandleBestSellerRequests(
	w http.ResponseWriter,
	r *http.Request,
	cat string,
	BestSellerFn book.BestSellerFn,
	pool *pgxpool.Pool,
	cache *ristretto.Cache[string, []byte],
) {
	libCode := r.URL.Query().Get("libCode")
	cacheKey := fmt.Sprintf("%s %s", cat, libCode)

	response, cacheFound := cache.Get(cacheKey)

	if cacheFound {
		log.Printf("hit, %v", cacheKey)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}

	if !cacheFound {
		ctx := context.Background()
		conn, err := pool.Acquire(ctx)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Release()
		query := sqlc.New(conn)

		resp := BestSellerFn()
		var isbns []string
		start := time.Now()
		for _, item := range resp.Items {
			isbns = append(isbns, item.Isbn13)
			query.UpdateDescription(ctx, sqlc.UpdateDescriptionParams{
				Description: pgtype.Text{String: item.Description, Valid: true},
				Isbn:        pgtype.Text{String: item.Isbn13, Valid: true},
				Source:      pgtype.Text{String: "aladin", Valid: true},
			})
			query.UpdateImageUrl(ctx, sqlc.UpdateImageUrlParams{
				ImageUrl: pgtype.Text{String: item.Cover, Valid: true},
				Isbn:     pgtype.Text{String: item.Isbn13, Valid: true},
				Source:   pgtype.Text{String: "aladin", Valid: true},
			})
		}
		end := time.Now()
		log.Printf("time : %vms", end.Sub(start).Milliseconds())

		filteredIsbns, err := query.ReturnExistIsbns(ctx, sqlc.ReturnExistIsbnsParams{
			Isbns:    isbns,
			LibCodes: strings.Split(libCode, ","),
		})

		if err != nil {
			log.Printf("err: %#+v\n", err)
			http.Error(w, "db data encoding error", http.StatusInternalServerError)
		}

		var result []book.AladinItem
		for _, row := range filteredIsbns {
			for _, item := range resp.Items {
				if item.Isbn13 == row.Isbn.String {
					item.LibCode = row.LibCode
					result = append(result, item)
					break
				}
			}
		}
		response, err = json.Marshal(book.AladinResponse{
			CatValue: resp.CatValue,
			CatName:  resp.CatName,
			Items:    result,
		})
		if err != nil {
			log.Printf("err: %#+v\n", err)
			http.Error(w, "marshal Error", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		cache.SetWithTTL(cacheKey, response, 1, time.Duration(timeUntilNextMonday()))

	}
}

func timeUntilNextMonday() time.Duration {
	now := time.Now()
	weekday := now.Weekday()
	// 월요일이 1이므로, 현재 요일에서 다음 월요일까지의 남은 일수를 계산
	daysUntilMonday := (time.Monday - weekday + 7) % 7
	if daysUntilMonday == 0 {
		// 오늘이 월요일인 경우, 다음 주 월요일까지 계산
		daysUntilMonday = 7
	}
	nextMonday := now.AddDate(0, 0, int(daysUntilMonday)).Truncate(24 * time.Hour)
	return nextMonday.Sub(now)
}
