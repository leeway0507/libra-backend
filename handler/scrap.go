package handler

import (
	"context"
	"encoding/json"
	"libra-backend/db/sqlc"
	"libra-backend/pkg/scrap"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

func GetScrapRouter(query *sqlc.Queries) *http.ServeMux {
	searchRouter := http.NewServeMux()
	searchRouter.HandleFunc("GET /{libCode}/{isbn}", func(w http.ResponseWriter, r *http.Request) {
		HandleScrap(w, r, query)
	})
	return searchRouter
}

func HandleScrap(w http.ResponseWriter, r *http.Request, query *sqlc.Queries) {
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

	ctx := context.Background()
	for _, d := range *data {
		ReqLibCode, err := strconv.Atoi(libCode)
		if err != nil {
			log.Println("fail to convert String to Int")
			break
		}
		if strings.Contains(engine.GetLibName(), d.LibName) {
			query.UpdateClassNum(ctx, sqlc.UpdateClassNumParams{
				ClassNum: pgtype.Text{String: d.BookCode, Valid: true},
				Isbn:     pgtype.Text{String: isbn, Valid: true},
				LibCode:  pgtype.Int4{Int32: int32(ReqLibCode), Valid: true},
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
