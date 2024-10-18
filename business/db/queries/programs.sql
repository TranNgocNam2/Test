-- name: CreateProgram :exec
INSERT INTO programs (id, name, start_date, end_date, created_by, description)
VALUES (sqlc.arg(id), sqlc.arg(name), sqlc.arg(start_date), sqlc.arg(end_date),
        sqlc.arg(created_by), sqlc.arg(description));

-- name: UpdateProgram :exec
UPDATE programs
SET name = sqlc.arg(name),
    start_date = sqlc.arg(start_date),
    end_date = sqlc.arg(end_date),
    updated_by = sqlc.arg(updated_by),
    description = sqlc.arg(description),
    updated_at = now()
WHERE id = sqlc.arg(id);

-- name: GetProgramByID :one
SELECT * FROM programs WHERE id = sqlc.arg(id);

-- name: DeleteProgram :exec
DELETE FROM programs WHERE id = sqlc.arg(id);