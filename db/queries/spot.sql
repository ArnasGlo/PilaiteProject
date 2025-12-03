-- name: InsertSpot :one
INSERT INTO spot(
    category, name, description, location_id
) VALUES (
             $1, $2, $3, $4
         ) RETURNING *;

-- name: GetSpotByID :one
SELECT * FROM spot WHERE id = $1;

-- name: GetAllSpots :many
SELECT * FROM spot;

-- name: GetPublicSpotsWithDetails :many
SELECT
    s.id,
    s.name,
    s.category,
    l.address,
    l.latitude,
    l.longitude,
    COALESCE((SELECT url FROM image WHERE spot_id = s.id LIMIT 1), '')::text as image_url
FROM spot s
         INNER JOIN location l ON s.location_id = l.id
WHERE s.category != 'Slaptos_vietos'
ORDER BY s.id;

-- name: GetSpotsWithDetails :many
SELECT
    s.id,
    s.name,
    s.category,
    l.address,
    l.latitude,
    l.longitude,
    COALESCE((SELECT url FROM image WHERE spot_id = s.id LIMIT 1), '')::text as image_url
FROM spot s
         INNER JOIN location l ON s.location_id = l.id
ORDER BY s.id;

-- name: GetPublicSpotsByCategoryWithDetails :many
SELECT
    s.id,
    s.name,
    s.category,
    l.address,
    l.latitude,
    l.longitude,
    COALESCE((SELECT url FROM image WHERE spot_id = s.id LIMIT 1), '')::text as image_url
FROM spot s
         INNER JOIN location l ON s.location_id = l.id
WHERE s.category = $1 AND s.category != 'Slaptos_vietos'
ORDER BY s.id;

-- name: GetSpotsByCategoryWithDetails :many
SELECT
    s.id,
    s.name,
    s.category,
    l.address,
    l.latitude,
    l.longitude,
    COALESCE((SELECT url FROM image WHERE spot_id = s.id LIMIT 1), '')::text AS image_url
FROM spot s
         INNER JOIN location l ON s.location_id = l.id
WHERE s.category = $1
ORDER BY s.id;
