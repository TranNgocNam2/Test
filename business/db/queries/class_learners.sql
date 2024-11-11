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

-- name: GetLearnersByClassId :many
SELECT u.*, cl.id AS class_learner_id, s.id AS school_id, s.name AS school_name
FROM users u
        JOIN class_learners cl ON cl.learner_id = u.id
        JOIN classes c ON cl.class_id = c.id
        JOIN schools s ON s.id = u.school_id
WHERE c.id = sqlc.arg(class_id)::uuid;