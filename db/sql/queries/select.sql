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
SELECT DISTINCT ON (b.isbn) 
	b.isbn,
	b.title,
	b.author,
	b.publisher,
	b.publication_Year,
	b.image_Url,
((bigm_similarity(author, @keyword) + bigm_similarity(title, @keyword)) * 10)::FLOAT AS score
FROM books b
JOIN libsbooks l
ON b.isbn = l.isbn
WHERE (b.author LIKE '%' || @keyword || '%' OR b.title LIKE '%' || @keyword || '%') 
        AND l.lib_code = ANY(@lib_codes::int[]) 
ORDER BY b.isbn DESC
LIMIT 50;

-- explain (analyze)
-- SELECT DISTINCT ON (b.isbn) 
--   b.*, 
--   ((bigm_similarity(b.author, '파이썬') + bigm_similarity(b.title, '파이썬')) * 10) AS score
-- FROM books b
-- JOIN libsbooks l
-- ON b.isbn = l.isbn
-- WHERE (b.author LIKE '%파이썬%' OR b.title LIKE '%파이썬%') 
--   AND l.lib_code = ANY(ARRAY[711539,111003,111006,111005,111004,111501])
-- ORDER BY b.isbn, score DESC
-- LIMIT 50;
