-- name: UpdateDescription :exec
UPDATE books
SET
    Description = $1
WHERE
    isbn = $2
    AND (
        Description = ''
        OR Description IS NULL
    );

-- name: UpdateRecom :exec
UPDATE books
SET
    Recommendation = $1
WHERE
    isbn = $2
    AND (
        Recommendation = ''
        OR Recommendation IS NULL
    );

-- name: UpdateToc :exec
UPDATE books
SET
    Toc = $1
WHERE
    isbn = $2
    AND (
        Toc = ''
        OR Toc IS NULL
    );

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