-- name: GetAssignmentsByClassLearner :many
SELECT * FROM learner_assignments
    WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid;