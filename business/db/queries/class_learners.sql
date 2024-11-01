-- name: CountLearnersByClassId :one
SELECT COUNT(*) FROM class_learners WHERE class_id = sqlc.arg(class_id)::uuid;