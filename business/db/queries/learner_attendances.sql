-- name: GenerateLearnerAttendance :exec
INSERT INTO learner_attendances(id, class_learner_id, slot_id)
VALUES (uuid_generate_v4(), sqlc.arg(class_learner_id), sqlc.arg(slot_id)::uuid);

-- name: SubmitLearnerAttendance :exec
UPDATE learner_attendances
SET status = sqlc.arg(status)
WHERE id = sqlc.arg(id)::uuid;

-- name: GetLearnerAttendanceByClassLearnerAndSlot :one
SELECT * FROM learner_attendances
    WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid
    AND slot_id = sqlc.arg(slot_id)::uuid;

-- name: GetAttendanceByClassLearner :many
SELECT * FROM learner_attendances
    WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid;

-- name: GetLearnerAttendanceBySlot :many
SELECT
    u.id, u.full_name, s.id AS school_id, s.name AS school_name, la.status
FROM users u
         JOIN class_learners cl ON u.id = cl.learner_id
         JOIN schools s ON s.id = u.school_id
         JOIN learner_attendances la ON la.class_learner_id = cl.id
WHERE la.slot_id = sqlc.arg(slot_id)::uuid;