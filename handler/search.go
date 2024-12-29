package handler

import (
	"context"
	"encoding/json"
	"libra-backend/db/sqlc"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetSearchRouter(query *sqlc.Queries) *http.ServeMux {
	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /normal", func(w http.ResponseWriter, r *http.Request) {
		HandleSearchNormal(w, r, query)
	})
	return searchRouter
}

func HandleSearchNormal(w http.ResponseWriter, r *http.Request, query *sqlc.Queries) {
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
	var libCodeIntArr []int32
	for _, code := range libCodeArr {
		intCode, err := strconv.Atoi(code)

		if err != nil {
			log.Printf("Error converting libCode to int: %v", err)
			http.Error(w, "failed to convert libCode to Int", http.StatusInternalServerError)
			return
		}

		libCodeIntArr = append(libCodeIntArr, int32(intCode))
	}

	ctx := context.Background()

	data, err := query.SearchFromBooks(ctx, sqlc.SearchFromBooksParams{
		Keyword:  pgtype.Text{String: keyword, Valid: true},
		LibCodes: libCodeIntArr,
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
}
