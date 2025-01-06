package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/scrap"
	"log"
	"net/http"
	"strings"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetScrapRouter(pool *pgxpool.Pool) (*http.ServeMux, func()) {
	searchRouter := http.NewServeMux()
	cache, err := ristretto.NewCache(&ristretto.Config[string, []byte]{
		NumCounters: 1e3,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 28, // maximum cost of cache (1GB).
		BufferItems: 8,       // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	searchRouter.HandleFunc("GET /{libCode}/{isbn}", func(w http.ResponseWriter, r *http.Request) {
		HandleScrap(w, r, pool, cache)
	})
	return searchRouter, cache.Close
}

func HandleScrap(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool, cache *ristretto.Cache[string, []byte]) {
	ctx := context.Background()
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Release()

	query := sqlc.New(conn)

	libCode := r.PathValue("libCode")
	isbn := r.PathValue("isbn")
	if isbn == "" {
		http.Error(w, "ISBN not provided", http.StatusBadRequest)
		return
	}
	// get value from cache
	cacheKey := fmt.Sprintf("%s %s", isbn, libCode)
	response, found := cache.Get(cacheKey)

	if found {
		log.Printf("hit, %v", cacheKey)
	}

	if !found {
		scraper := scrap.GetInstance(libCode)

		if scraper == nil {
			log.Printf("Not Found Lib scraper : %s", libCode)
			http.Error(w, "Not Found Lib scraper", http.StatusInternalServerError)
			return
		}

		scraperInstance := scraper(isbn)

		log.Printf("request info \n %#+v\n libCode:%s\n isbn:%s\n", scraperInstance.GetDistrict(), libCode, isbn)

		reader, err := scraperInstance.Request()
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data, err := scraperInstance.ExtractData(reader)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Not found info matched with ISBN", http.StatusNotFound)
			return
		}

		for _, d := range *data {
			if strings.Contains(scraperInstance.GetLibName(), d.LibName) {
				query.UpdateClassNum(ctx, sqlc.UpdateClassNumParams{
					ClassNum: pgtype.Text{String: d.BookCode, Valid: true},
					Isbn:     pgtype.Text{String: isbn, Valid: true},
					LibCode:  pgtype.Text{String: libCode, Valid: true},
				})
				break
			}
		}

		response, err = json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		cache.Set(cacheKey, response, 1)
		cache.Wait()
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
