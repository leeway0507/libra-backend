-- name: GerSearchResult :many
SELECT * FROM Books;

-- name: GetBookDetail :one
SELECT 
    b.*,
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'libCode', l.lib_code,
            'classNum', l.class_num,
            'bookCode', l.book_code
        )
    ) AS lib_books
FROM Books b
JOIN libsbooks l 
    ON b.isbn = l.isbn AND l.lib_code = ANY(@lib_codes::int[])
WHERE b.isbn = @isbn
GROUP BY b.isbn, b.id;



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
SELECT * 
FROM books
WHERE author LIKE '%' || @keyword || '%' OR title LIKE '%' || @keyword || '%'
ORDER BY ((bigm_similarity(author, @keyword) + bigm_similarity(title, @keyword)) * 10) DESC
LIMIT 50;