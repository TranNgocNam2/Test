-- name: GetDistrictsByProvince :many
SELECT * FROM districts
         WHERE districts.province_id = sqlc.arg(province_id)::integer
         ORDER BY id;

-- name: GetDistrictById :one
SELECT * FROM districts
         WHERE districts.id = sqlc.arg(id)::integer;