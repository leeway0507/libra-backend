-- name: GetSearchResult :many
SELECT * FROM Books;

-- name: GetBookDetail :one
SELECT b.*, JSON_AGG(
        JSON_BUILD_OBJECT(
            'libCode', l.lib_code, 'classNum', l.class_num, 'bookCode', l.book_code
        )
    ) AS lib_books
FROM Books b
    JOIN libsbooks l ON b.isbn = l.isbn
    AND l.lib_code = ANY (@ lib_codes::VARCHAR(20) [])
WHERE
    b.isbn = @ isbn
GROUP BY
    b.isbn,
    b.id;

-- name: ExtractBooksForEmbedding :many
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
    );
-- name: GetLibCodFromLibName :one
SELECT lib_code FROM libraries WHERE lib_name = $1;

-- name: SearchFromBooks :many
WITH
    FilteredLibsBooks AS (
        SELECT DISTINCT
            isbn
        FROM libsbooks
        WHERE
            lib_code = ANY (@ lib_codes::VARCHAR(15) [])
    )
SELECT b.isbn, b.title, b.author, b.publisher, b.publication_year, b.image_url, (embedding <= > $1)::REAL as score
FROM
    books b
    JOIN BookEmbedding e ON b.isbn = e.isbn
    JOIN FilteredLibsBooks l ON l.isbn = e.isbn
WHERE
    embedding <= > $1 <= 0.8
ORDER BY embedding <= > $1 ASC
LIMIT 50;