-- name: InsertAssignment :one
INSERT INTO assignments (id, classId, question, deadline, status, type, can_overdue)
VALUES (sqlc.arg(id)::uuid, sqlc.arg(classId)::uuid, sqlc.arg(question), sqlc.arg(deadline),
        sqlc.arg(status), sqlc.arg(type), sqlc.arg(can_overdue))
RETURNING id;

-- name: InsertLearnerAssignment :copyfrom
INSERT INTO learner_assignments(id, class_learner_id, assignment_id, grade, grading_status, submission_status)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateAssignment :exec
UPDATE assignments
SET question = sqlc.arg(question),
    deadline = sqlc.arg(deadline),
    status = sqlc.arg(status),
    type = sqlc.arg(type),
    can_overdue = sqlc.arg(can_overdue)
WHERE id = sqlc.arg(id)::uuid;

-- name: CheckAssignmentInClass :one
SELECT EXISTS (SELECT 1
FROM assignments WHERE
    class_id = sqlc.arg(class_id)::uuid
AND id = sqlc.arg(id)::uuid);

-- name: GetAssignmentById :one
SELECT * FROM assignments WHERE id = sqlc.arg(id)::uuid;

-- name: DeleteAssignment :exec
DELETE FROM assignments WHERE id = sqlc.arg(id)::uuid;

-- name: UpdateLearnerGrade :exec
UPDATE learner_assignments
SET grading_status = sqlc.arg(grading_status),
    grade = sqlc.arg(grade)
WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid
AND assignment_id = sqlc.arg(assignment_id)::uuid;

-- name: GetLearnerAssignment :one
SELECT * from learner_assignments Where class_learner_id = sqlc.arg(class_learner_id) AND assignment_id = sqlc.arg(assignment_id);

-- name: UpdateLearnerAssignment :exec
UPDATE learner_assignments
SET data = sqlc.arg(data),
    submission_status = sqlc.arg(submission_status)
WHERE class_learner_id = sqlc.arg(class_learner_id)::uuid
AND assignment_id = sqlc.arg(assignment_id)::uuid;
