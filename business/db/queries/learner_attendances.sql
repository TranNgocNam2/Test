-- name: GenerateLearnerAttendance :exec
INSERT INTO learner_attendances(id, class_learner_id, slot_id)
VALUES (uuid_generate_v4(), sqlc.arg(class_learner_id), UNNEST(sqlc.arg(slot_ids)::uuid[]));

-- name: GenerateLearnersAttendance :exec
INSERT INTO learner_attendances(id, class_learner_id, slot_id)
SELECT
    uuid_generate_v4(),
    class_learner_id,
    slot_id
FROM
    UNNEST(sqlc.arg(class_learner_ids)::uuid[]) AS learner_ids(class_learner_id)
        CROSS JOIN
    UNNEST(sqlc.arg(slot_ids)::uuid[]) AS slot_ids(slot_id);

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
         JOIN verification_learners vl ON u.id = vl.learner_id
         JOIN schools s ON s.id = vl.school_id
         JOIN learner_attendances la ON la.class_learner_id = cl.id
WHERE la.slot_id = sqlc.arg(slot_id)::uuid;

-- name: GetLearnerAttendanceRecords :many
SELECT la.id, la.status, s.start_time, s.end_time, s.index
FROM learner_attendances la
         JOIN slots s ON s.id = la.slot_id
         JOIN class_learners cl ON cl.id = la.class_learner_id
WHERE cl.class_id = sqlc.arg(class_id)::uuid
  AND cl.learner_id = sqlc.arg(learner_id)
ORDER BY s.index;