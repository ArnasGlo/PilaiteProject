-- name: InsertImage :one
INSERT INTO image(
    url, spot_id
) VALUES (
             $1, $2
         )RETURNING *;

-- name: GetImageByID :one
SELECT * FROM image WHERE id = $1;

-- name: GetAllImages :many
SELECT * FROM image;