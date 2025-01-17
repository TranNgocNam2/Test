-- name: CreateSchool :exec
INSERT INTO schools (id, name, address, district_id)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(address), sqlc.arg(district_id)::integer);

-- name: DeleteSchool :exec
UPDATE schools
SET is_deleted = true
WHERE id = sqlc.arg(id)::uuid;

-- name: GetSchoolById :one
SELECT * FROM schools
WHERE id = sqlc.arg(id)::uuid AND is_deleted = false;

-- name: UpdateSchool :exec
UPDATE schools
SET name = sqlc.arg(name), address = sqlc.arg(address), district_id = sqlc.arg(district_id)::integer
WHERE id = sqlc.arg(id)::uuid;

-- name: GetSchoolsByDistrictId :many
SELECT * FROM schools
WHERE district_id = sqlc.arg(district_id)::integer
AND is_deleted = false;

-- name: GetAllSchools :many
SELECT s.*, d.id AS district_id, d.name AS district_name, p.id AS province_id, p.name AS province_name
FROM schools s JOIN districts d ON s.district_id = d.id
               JOIN provinces p ON p.id = d.province_id;

-- name: GetAllSchoolInformationById :one
SELECT s.*, d.id AS district_id, d.name AS district_name, p.id AS province_id, p.name AS province_name
FROM schools s JOIN districts d ON s.district_id = d.id
               JOIN provinces p ON p.id = d.province_id
WHERE s.id = sqlc.arg(id)::uuid;