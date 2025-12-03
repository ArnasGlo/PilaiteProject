-- name: InsertUser :one
INSERT INTO users(
    email, password, role
) VALUES (
             $1, $2, $3
         )RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

