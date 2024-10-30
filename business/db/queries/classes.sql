-- name: CreateClass :exec
INSERT INTO classes (id, code, password, name, subject_id, program_id, link, start_date, end_date, created_by)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(code), sqlc.arg(password),
        sqlc.arg(name), sqlc.arg(subject_id)::uuid, sqlc.arg(program_id)::uuid,
        sqlc.arg(link), sqlc.arg(start_date), sqlc.arg(end_date), sqlc.arg(created_by));

-- name: GetClassById :one
SELECT * FROM classes WHERE id = sqlc.arg(id)::uuid;

-- name: CountClassesByProgramId :one
SELECT COUNT(*) FROM classes WHERE program_id = sqlc.arg(program_id)::uuid;

-- name: GetClassByCode :one
SELECT * FROM classes WHERE code = sqlc.arg(code);

-- name: UpdateActiveClass :exec
UPDATE classes
SET status = 1,
    start_date = sqlc.arg(start_date),
    end_date = sqlc.arg(end_date)
WHERE id = sqlc.arg(id)::uuid;

-- name: DeleteClass :exec
DELETE FROM classes WHERE id = sqlc.arg(id)::uuid;

-- name: SoftDeleteClass :exec
UPDATE classes
SET status = 2
WHERE id = sqlc.arg(id)::uuid;

-- name: UpdateClass :exec
UPDATE classes
SET name = sqlc.arg(name),
    code = sqlc.arg(code),
    password = sqlc.arg(password)
WHERE id = sqlc.arg(id)::uuid;