-- name: UpdateDescription :exec
UPDATE books
SET
    Description = $1,
    source = $2
WHERE
    isbn = $3
    AND source != $2;
-- name: UpdateImageUrl :exec
UPDATE books
SET
    image_url = $1,
    source = $2
WHERE
    isbn = $3
    AND (
        image_url = ''
        OR image_url IS NULL
    );

-- name: UpdateRecom :exec
UPDATE books
SET
    Recommendation = $1,
    source = $2
WHERE
    isbn = $3
    AND source != $2;

-- name: UpdateToc :exec
UPDATE books
SET
    Toc = $1,
    source = $2
WHERE
    isbn = $3
    AND source != $2;

-- name: UpdateClassNum :exec
UPDATE libsbooks
SET
    class_num = $1,
    scrap = true
WHERE
    isbn = $2
    and lib_code = $3
    AND (
        class_num = ''
        OR class_num IS NULL
    );