-- name: CreateSpecialization :exec
INSERT INTO specializations (id, name, code, time_amount, image_link, description, created_by)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.arg(code), sqlc.arg(time_amount),
        sqlc.arg(image_link), sqlc.arg(description), sqlc.arg(created_by));

