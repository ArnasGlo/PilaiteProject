-- name: InsertLocation :one
INSERT INTO location(
    address, latitude, longitude
) VALUES (
             $1, $2, $3
         )RETURNING *;

-- name: GetLocationByID :one
SELECT * FROM location WHERE id = $1;

-- name: GetAllLocations :many
SELECT * FROM location;