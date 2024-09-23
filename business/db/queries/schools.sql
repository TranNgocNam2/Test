-- name: CreateSchool :exec
INSERT INTO schools (id, name, address, district_id)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(address), sqlc.arg(district_id)::integer);

-- name: DeleteSchool :exec
DELETE FROM schools
WHERE id = sqlc.arg(id)::uuid;

-- name: GetSchoolByID :one
SELECT * FROM schools
WHERE id = sqlc.arg(id)::uuid;