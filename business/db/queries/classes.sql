-- name: CreateClass :exec
INSERT INTO classes (id, code, password, name, subject_id, program_id, link, start_date, end_date, created_by)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(code), sqlc.arg(password),
        sqlc.arg(name), sqlc.arg(subject_id)::uuid, sqlc.arg(program_id)::uuid,
        sqlc.arg(link), sqlc.arg(start_date), sqlc.arg(end_date), sqlc.arg(created_by));

-- name: GetClassByID :one
SELECT * FROM classes WHERE id = sqlc.arg(id)::uuid;

-- name: CountClassesByProgramID :one
SELECT COUNT(*) FROM classes WHERE program_id = sqlc.arg(program_id)::uuid;