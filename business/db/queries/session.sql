-- name: CountSessionsBySubjectID :one
SELECT count(*) FROM sessions WHERE subject_id = sqlc.arg(subject_id);