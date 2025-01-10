package book

import (
	"context"
	"encoding/json"
	"libra-backend/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type LibArr struct {
	LibCode  string `json:"libCode"`
	Classnum string `json:"classNum"`
}
type BookDetail struct {
	sqlc.GetBookDetailRow
	LibBooks []LibArr `json:"libBooks"`
}

func RequestBookDetail(query *sqlc.Queries, isbn string, libCodes []string) (*BookDetail, error) {
	params := sqlc.GetBookDetailParams{
		Isbn:     pgtype.Text{String: isbn, Valid: true},
		LibCodes: libCodes,
	}
	ctx := context.Background()
	detail, err := query.GetBookDetail(ctx, params)
	if err != nil {
		return nil, err
	}
	var LibBooks []LibArr
	err = json.Unmarshal(detail.LibBooks, &LibBooks)
	if err != nil {
		return nil, err
	}

	return &BookDetail{detail, LibBooks}, nil

}
