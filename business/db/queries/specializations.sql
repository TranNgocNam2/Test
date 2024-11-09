-- name: CreateSpecialization :exec
INSERT INTO specializations (id, name, code, time_amount, image_link, description, created_by)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.arg(code), sqlc.arg(time_amount),
        sqlc.arg(image_link), sqlc.arg(description), sqlc.arg(created_by));

-- name: GetSpecializationByCode :one
SELECT * FROM specializations WHERE code = sqlc.arg(code);

-- name: GetSpecializationById :one
SELECT * FROM specializations WHERE id = sqlc.arg(id);

-- name: UpdateSpecialization :exec
UPDATE specializations SET name = sqlc.arg(name), code = sqlc.arg(code), status = sqlc.arg(status), time_amount = sqlc.arg(time_amount),
        image_link = sqlc.arg(image_link), description = sqlc.arg(description), updated_at = NOW(), updated_by = sqlc.arg(updated_by)
WHERE id = sqlc.arg(id);

-- name: UpdateSpecializationStatus :exec
UPDATE specializations SET status = 2, updated_at = NOW(), updated_by = sqlc.arg(updated_by)
WHERE id = sqlc.arg(id);

-- name: DeleteSpecialization :exec
DELETE FROM specializations WHERE id = sqlc.arg(id) AND status = 0;

-- name: GetPublishedSpecializationById :one
SELECT * FROM specializations WHERE id = sqlc.arg(id) AND status = 1;