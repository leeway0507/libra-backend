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
