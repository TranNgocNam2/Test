-- name: InsertSubject :one
INSERT INTO subjects (id, name, code, description, image_link, status,
    time_per_session, sessions_per_week, created_by,
    created_at)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(name), sqlc.arg(code), sqlc.arg(description),
    sqlc.arg(image_link), sqlc.arg(status), sqlc.arg(time_per_session),
    sqlc.arg(sessions_per_week), sqlc.arg(created_by),
    sqlc.arg(created_at))
RETURNING id;

-- name: InsertSubjectSkill :copyfrom
INSERT INTO subject_skills (subject_id, skill_id)
VALUES ($1, $2);
