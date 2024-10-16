-- name: UpsertSession :exec
INSERT INTO sessions(id, subject_id, index, name)
VALUES(sqlc.arg(id)::uuid, sqlc.arg(subject_id)::uuid, sqlc.arg(index), sqlc.arg(name))
ON CONFLICT (id)
DO UPDATE SET
    index = EXCLUDED.index,
    name = EXCLUDED.name;


-- name: CountSessionsBySubjectID :one
SELECT count(*) FROM sessions WHERE subject_id = sqlc.arg(subject_id);

-- name: GetSessionsBySubjectID :many
SELECT * FROM sessions WHERE subject_id = sqlc.arg(subject_id) ORDER BY index;
