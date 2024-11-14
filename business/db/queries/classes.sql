-- name: CreateClass :exec
INSERT INTO classes (id, code, password, name, subject_id, program_id, link, start_date, end_date, created_by)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(code), sqlc.arg(password),
        sqlc.arg(name), sqlc.arg(subject_id)::uuid, sqlc.arg(program_id)::uuid,
        sqlc.arg(link), sqlc.arg(start_date), sqlc.arg(end_date), sqlc.arg(created_by));

-- name: GetClassById :one
SELECT * FROM classes WHERE id = sqlc.arg(id)::uuid;

-- name: CountClassesByProgramId :one
SELECT COUNT(*) FROM classes WHERE program_id = sqlc.arg(program_id)::uuid;

-- name: GetClassCompletedByCode :one
SELECT * FROM classes WHERE code = sqlc.arg(code) AND status = 1;

-- name: UpdateClassStatusAndDate :exec
UPDATE classes
SET status = sqlc.arg(status),
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

-- name: CheckTeacherInClass :one
SELECT EXISTS (SELECT 1
FROM slots WHERE
    teacher_id = sqlc.arg(teacher_id)
AND class_id = sqlc.arg(class_id)::uuid);

-- name: UpdateMeetingLink :exec
UPDATE classes
SET link = sqlc.arg(link),
    updated_at = now(),
    updated_by = sqlc.arg(updated_by)
WHERE id = sqlc.arg(id)::uuid;