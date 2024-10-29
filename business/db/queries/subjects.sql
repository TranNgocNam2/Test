-- name: InsertSubject :one
INSERT INTO subjects (id, name, code, description, image_link, status,
    time_per_session, sessions_per_week, created_by,
    created_at)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(code), sqlc.arg(description),
    sqlc.arg(image_link), sqlc.arg(status), sqlc.arg(time_per_session),
    sqlc.arg(sessions_per_week), sqlc.arg(created_by),
    sqlc.arg(created_at))
RETURNING id;

-- name: DeleteSubjectSkills :exec
DELETE FROM subject_skills WHERE subject_id = sqlc.arg(subject_id);

-- name: GetSubjectsByIDs :many
SELECT *
FROM subjects
WHERE id = ANY(sqlc.arg(subject_ids)::uuid[]) AND status = 1;

-- name: GetSubjectByCode :one
SELECT *
FROM subjects
WHERE code = sqlc.arg(code);

-- name: UpdateSubject :exec
UPDATE subjects
SET name = sqlc.arg(name),
    code = sqlc.arg(code),
    description = sqlc.arg(description),
    status = sqlc.arg(status),
    image_link = sqlc.arg(image_link),
    updated_by = sqlc.arg(updated_by),
    updated_at = sqlc.arg(updated_at)
WHERE id = sqlc.arg(id)::uuid;

-- name: GetSubjectById :one
SELECT *
FROM subjects WHERE id = sqlc.arg(id)::uuid;

-- name: DeleteSubject :exec
UPDATE subjects SET status = 2, updated_at = NOW(), updated_by = sqlc.arg(updated_by)
WHERE id = sqlc.arg(id);

-- name: GetPublishedSubjectByID :one
SELECT *
FROM subjects
WHERE id = sqlc.arg(id)::uuid AND status = 1;

-- -- name: GetSubjectsByIDs :many
-- SELECT * FROM subjects WHERE id = ANY(sqlc.arg(subject_ids)::uuid[]);