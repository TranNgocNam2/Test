-- name: DeleteSubjectTranscripts :exec
DELETE FROM transcripts WHERE subject_id = sqlc.arg(subject_id);

-- name: InsertTranscripts :copyfrom
INSERT INTO transcripts(id, subject_id, name, index, min_grade, weight)
VALUES($1, $2, $3, $4, $5, $6);

-- name: GetTranscriptsCountBySubjectId :one
SELECT COUNT(*), COALESCE(SUM(weight), 0)::integer AS sum  FROM transcripts WHERE subject_id = sqlc.arg(subject_id);

-- name: GetTranscriptsBySubjectId :many
SELECT * FROM transcripts WHERE subject_id = sqlc.arg(subject_id) ORDER BY index;


