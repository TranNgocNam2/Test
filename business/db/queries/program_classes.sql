-- name: CountClassesByProgramID :one
SELECT COUNT(*) FROM program_classes
WHERE program_id = sqlc.arg(program_id)::uuid;