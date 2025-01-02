package handler

import (
	"context"
	"encoding/json"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/scrap"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetScrapRouter(pool *pgxpool.Pool) *http.ServeMux {
	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /{libCode}/{isbn}", func(w http.ResponseWriter, r *http.Request) {
		HandleScrap(w, r, pool)
	})
	return searchRouter
}

func HandleScrap(w http.ResponseWriter, r *http.Request, pool *pgxpool.Pool) {
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

	engine := scrap.GetInstance(libCode)(isbn)

	log.Printf("request info \n %#+v\n libCode:%s\n isbn:%s\n", engine.GetDistrict(), libCode, isbn)

	reader, err := engine.Request()
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data, err := engine.ExtractData(reader)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, "Not found info matched with ISBN", http.StatusNotFound)
		return
	}

	for _, d := range *data {

		if strings.Contains(engine.GetLibName(), d.LibName) {
			query.UpdateClassNum(ctx, sqlc.UpdateClassNumParams{
				ClassNum: pgtype.Text{String: d.BookCode, Valid: true},
				Isbn:     pgtype.Text{String: isbn, Valid: true},
				LibCode:  pgtype.Text{String: libCode, Valid: true},
			})
			break
		}
		// query.UpdateDescription(ctx,sqlc.UpdateDescriptionParams{
		// 	Description: pgtype.Text{String: d.},
		// })
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
