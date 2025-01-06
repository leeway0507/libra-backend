// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

type Book struct {
	ID              int32       `json:"id"`
	Isbn            pgtype.Text `json:"isbn"`
	Title           pgtype.Text `json:"title"`
	Author          pgtype.Text `json:"author"`
	Publisher       pgtype.Text `json:"publisher"`
	PublicationYear pgtype.Text `json:"publicationYear"`
	Volume          pgtype.Text `json:"volume"`
	ImageUrl        pgtype.Text `json:"imageUrl"`
	Description     pgtype.Text `json:"description"`
	Recommendation  pgtype.Text `json:"recommendation"`
	Toc             pgtype.Text `json:"toc"`
	Source          pgtype.Text `json:"source"`
	Url             pgtype.Text `json:"url"`
	VectorSearch    pgtype.Bool `json:"vectorSearch"`
}

type Bookembedding struct {
	ID        int32           `json:"id"`
	Isbn      string          `json:"isbn"`
	Embedding pgvector.Vector `json:"embedding"`
}

type Library struct {
	ID            int32         `json:"id"`
	LibCode       pgtype.Text   `json:"libCode"`
	LibName       pgtype.Text   `json:"libName"`
	Address       pgtype.Text   `json:"address"`
	Tel           pgtype.Text   `json:"tel"`
	Latitude      pgtype.Float8 `json:"latitude"`
	Longitude     pgtype.Float8 `json:"longitude"`
	Homepage      pgtype.Text   `json:"homepage"`
	Closed        pgtype.Text   `json:"closed"`
	OperatingTime pgtype.Text   `json:"operatingTime"`
}

type Libsbook struct {
	ID       int32       `json:"id"`
	LibCode  pgtype.Text `json:"libCode"`
	Isbn     pgtype.Text `json:"isbn"`
	ClassNum pgtype.Text `json:"classNum"`
	Scrap    pgtype.Bool `json:"scrap"`
}
