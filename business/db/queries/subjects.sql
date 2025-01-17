-- name: InsertSubject :one
INSERT INTO subjects (id, name, code, description, image_link, status, sessions_per_week,
    time_per_session, created_by, created_at, learner_type)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(code), sqlc.arg(description),
        sqlc.arg(image_link), sqlc.arg(status), sqlc.arg(sessions_per_week),
        sqlc.arg(time_per_session), sqlc.arg(created_by), sqlc.arg(created_at),
        sqlc.arg(learner_type))
RETURNING id;

-- name: DeleteSubjectSkills :exec
DELETE FROM subject_skills WHERE subject_id = sqlc.arg(subject_id);

-- name: GetSubjectsByIds :many
SELECT *
FROM subjects
WHERE id = ANY(sqlc.arg(subject_ids)::uuid[]) AND status = 1;

-- name: IsSubjectCodePublished :one
SELECT *
FROM subjects
WHERE code = sqlc.arg(code) AND status = 1 AND id != sqlc.arg(id);

-- name: IsSubjectCodeExist :one
SELECT *
FROM subjects
WHERE code = sqlc.arg(code) AND status = 1;

-- name: UpdateSubject :exec
UPDATE subjects
SET name = sqlc.arg(name),
    code = sqlc.arg(code),
    time_per_session = sqlc.arg(time_per_session),
    sessions_per_week = sqlc.arg(sessions_per_week),
    total_sessions = sqlc.arg(total_sessions),
    min_pass_grade = sqlc.arg(min_pass_grade),
    min_attendance = sqlc.arg(min_attendance),
    description = sqlc.arg(description),
    status = sqlc.arg(status),
    image_link = sqlc.arg(image_link),
    updated_by = sqlc.arg(updated_by),
    updated_at = sqlc.arg(updated_at),
    learner_type = sqlc.arg(learner_type)
WHERE id = sqlc.arg(id)::uuid;

-- name: GetSubjectById :one
SELECT *
FROM subjects WHERE id = sqlc.arg(id)::uuid;

-- name: DeleteSubject :exec
UPDATE subjects SET status = 2, updated_at = NOW(), updated_by = sqlc.arg(updated_by)
WHERE id = sqlc.arg(id);


-- name: GetPublishedSubjectById :one
SELECT *
FROM subjects
WHERE id = sqlc.arg(id)::uuid AND status = 1;
