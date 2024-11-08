-- name: CountLearnersByClassId :one
SELECT COUNT(*) FROM class_learners WHERE class_id = sqlc.arg(class_id)::uuid;

-- name: AddLearnerToClass :exec
INSERT INTO class_learners (id, class_id, learner_id)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(class_id)::uuid, sqlc.arg(learner_id));

-- name: GetLearnerByClassId :one
SELECT * FROM class_learners
         WHERE class_id = sqlc.arg(class_id)::uuid
           AND learner_id = sqlc.arg(learner_id);

-- name: GetClassesByLearnerId :many
SELECT * FROM classes
WHERE id IN (SELECT class_id FROM class_learners WHERE learner_id = sqlc.arg(learner_id));