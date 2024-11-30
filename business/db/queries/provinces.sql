-- name: GetAllProvince :many
SELECT * FROM provinces
         ORDER BY id;

-- name: GetProvinceById :one
SELECT * FROM provinces
         WHERE provinces.id = sqlc.arg(id)::integer;