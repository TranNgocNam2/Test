-- name: DeleteSubjectTranscripts :exec
DELETE FROM transcripts WHERE subject_id = sqlc.arg(subject_id);

-- name: InsertTranscripts :copyfrom
INSERT INTO transcripts(id, subject_id, name, index, min_grade, weight)
VALUES($1, $2, $3, $4, $5, $6);

-- name: GetTranscriptsBySubjectId :many
SELECT * FROM transcripts WHERE subject_id = sqlc.arg(subject_id) ORDER BY index;

-- name: GetTranscriptIdsBySubjectId :many
SELECT id AS ids FROM transcripts
    WHERE subject_id = sqlc.arg(subject_id) ORDER BY index;
