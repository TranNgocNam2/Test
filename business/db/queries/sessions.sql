-- name: UpsertSession :exec
INSERT INTO sessions(id, subject_id, index, name)
VALUES(sqlc.arg(id)::uuid, sqlc.arg(subject_id)::uuid, sqlc.arg(index), sqlc.arg(name))
ON CONFLICT (id)
DO UPDATE SET
    index = EXCLUDED.index,
    name = EXCLUDED.name;


-- name: CountSessionsBySubjectId :one
SELECT count(*) FROM sessions WHERE subject_id = sqlc.arg(subject_id);

-- name: GetSessionsBySubjectId :many
SELECT * FROM sessions WHERE subject_id = sqlc.arg(subject_id) ORDER BY index;

-- name: GetSessionById :one
SELECT * FROM sessions WHERE id = sqlc.arg(id);