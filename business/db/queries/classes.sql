-- name: CreateClass :exec
INSERT INTO classes (id, code, password, name, link, program_subject_id,
                     start_date, end_date, created_by)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(code), sqlc.arg(password),
        sqlc.arg(name), sqlc.arg(link), sqlc.arg(program_subject_id)::uuid,
        sqlc.arg(start_date), sqlc.arg(end_date), sqlc.arg(created_by));

-- name: GetClassByID :one
SELECT * FROM classes WHERE id = sqlc.arg(id)::uuid;