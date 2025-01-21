// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: select.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pgvector/pgvector-go"
)

const extractBooksForEmbedding = `-- name: ExtractBooksForEmbedding :many
SELECT
    isbn,
    title,
    description,
    toc,
    recommendation
FROM books b
WHERE (
        b.vector_search is false
        and b.source is not null
    )
`

type ExtractBooksForEmbeddingRow struct {
	Isbn           pgtype.Text `json:"isbn"`
	Title          pgtype.Text `json:"title"`
	Description    pgtype.Text `json:"description"`
	Toc            pgtype.Text `json:"toc"`
	Recommendation pgtype.Text `json:"recommendation"`
}

func (q *Queries) ExtractBooksForEmbedding(ctx context.Context) ([]ExtractBooksForEmbeddingRow, error) {
	rows, err := q.db.Query(ctx, extractBooksForEmbedding)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExtractBooksForEmbeddingRow
	for rows.Next() {
		var i ExtractBooksForEmbeddingRow
		if err := rows.Scan(
			&i.Isbn,
			&i.Title,
			&i.Description,
			&i.Toc,
			&i.Recommendation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBookDetail = `-- name: GetBookDetail :one
SELECT b.id, b.isbn, b.title, b.author, b.publisher, b.publication_year, b.volume, b.image_url, b.description, b.recommendation, b.toc, b.source, b.url, b.vector_search, JSON_AGG(
        JSON_BUILD_OBJECT(
            'libCode', l.lib_code, 'classNum', l.class_num
        )
    ) AS lib_books
FROM Books b
    JOIN libsbooks l ON b.isbn = l.isbn
    AND l.lib_code = ANY ($1::VARCHAR(20)[])
WHERE
    b.isbn = $2
GROUP BY
    b.isbn,
    b.id
`

type GetBookDetailParams struct {
	LibCodes []string    `json:"libCodes"`
	Isbn     pgtype.Text `json:"isbn"`
}

type GetBookDetailRow struct {
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
	LibBooks        []byte      `json:"libBooks"`
}

func (q *Queries) GetBookDetail(ctx context.Context, arg GetBookDetailParams) (GetBookDetailRow, error) {
	row := q.db.QueryRow(ctx, getBookDetail, arg.LibCodes, arg.Isbn)
	var i GetBookDetailRow
	err := row.Scan(
		&i.ID,
		&i.Isbn,
		&i.Title,
		&i.Author,
		&i.Publisher,
		&i.PublicationYear,
		&i.Volume,
		&i.ImageUrl,
		&i.Description,
		&i.Recommendation,
		&i.Toc,
		&i.Source,
		&i.Url,
		&i.VectorSearch,
		&i.LibBooks,
	)
	return i, err
}

const getLibCodFromLibName = `-- name: GetLibCodFromLibName :one
SELECT lib_code FROM libraries WHERE lib_name = $1
`

func (q *Queries) GetLibCodFromLibName(ctx context.Context, libName pgtype.Text) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, getLibCodFromLibName, libName)
	var lib_code pgtype.Text
	err := row.Scan(&lib_code)
	return lib_code, err
}

const getSearchResult = `-- name: GetSearchResult :many
SELECT id, isbn, title, author, publisher, publication_year, volume, image_url, description, recommendation, toc, source, url, vector_search FROM Books
`

func (q *Queries) GetSearchResult(ctx context.Context) ([]Book, error) {
	rows, err := q.db.Query(ctx, getSearchResult)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.ID,
			&i.Isbn,
			&i.Title,
			&i.Author,
			&i.Publisher,
			&i.PublicationYear,
			&i.Volume,
			&i.ImageUrl,
			&i.Description,
			&i.Recommendation,
			&i.Toc,
			&i.Source,
			&i.Url,
			&i.VectorSearch,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const returnExistIsbns = `-- name: ReturnExistIsbns :many
SELECT isbn, ARRAY_AGG(lib_code)::VARCHAR[] as lib_code
FROM libsbooks
WHERE
    isbn = ANY ($1::VARCHAR[])
    AND lib_code = ANY ($2::VARCHAR[])
GROUP BY isbn
`

type ReturnExistIsbnsParams struct {
	Isbns    []string `json:"isbns"`
	LibCodes []string `json:"libCodes"`
}

type ReturnExistIsbnsRow struct {
	Isbn    pgtype.Text `json:"isbn"`
	LibCode []string    `json:"libCode"`
}

func (q *Queries) ReturnExistIsbns(ctx context.Context, arg ReturnExistIsbnsParams) ([]ReturnExistIsbnsRow, error) {
	rows, err := q.db.Query(ctx, returnExistIsbns, arg.Isbns, arg.LibCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReturnExistIsbnsRow
	for rows.Next() {
		var i ReturnExistIsbnsRow
		if err := rows.Scan(&i.Isbn, &i.LibCode); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchFromBooks = `-- name: SearchFromBooks :many
WITH
    FilteredLibsBooks AS (
        SELECT DISTINCT
            isbn
        FROM libsbooks
        WHERE
            lib_code = ANY ($3::VARCHAR(15)[])
    )
SELECT b.isbn, b.title, b.author, b.publisher, b.publication_year, b.image_url, (embedding <=> $1)::REAL as score
FROM
    books b
    JOIN BookEmbedding e ON b.isbn = e.isbn
    JOIN FilteredLibsBooks l ON l.isbn = e.isbn
WHERE
    b.title LIKE '%' || $2 || '%' AND
    embedding <=> $1 <= 0.8
ORDER BY embedding <=> $1 ASC
LIMIT 50
`

type SearchFromBooksParams struct {
	Embedding pgvector.Vector `json:"embedding"`
	Keyword   pgtype.Text     `json:"keyword"`
	LibCodes  []string        `json:"libCodes"`
}

type SearchFromBooksRow struct {
	Isbn            pgtype.Text `json:"isbn"`
	Title           pgtype.Text `json:"title"`
	Author          pgtype.Text `json:"author"`
	Publisher       pgtype.Text `json:"publisher"`
	PublicationYear pgtype.Text `json:"publicationYear"`
	ImageUrl        pgtype.Text `json:"imageUrl"`
	Score           float32     `json:"score"`
}

func (q *Queries) SearchFromBooks(ctx context.Context, arg SearchFromBooksParams) ([]SearchFromBooksRow, error) {
	rows, err := q.db.Query(ctx, searchFromBooks, arg.Embedding, arg.Keyword, arg.LibCodes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SearchFromBooksRow
	for rows.Next() {
		var i SearchFromBooksRow
		if err := rows.Scan(
			&i.Isbn,
			&i.Title,
			&i.Author,
			&i.Publisher,
			&i.PublicationYear,
			&i.ImageUrl,
			&i.Score,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
