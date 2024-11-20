-- name: GetUserById :one
SELECT * FROM users
WHERE id = sqlc.arg(id);

-- name: CreateUser :exec
INSERT INTO users (id, email, auth_role, full_name)
VALUES (sqlc.arg(id), sqlc.arg(email), sqlc.arg(auth_role), sqlc.arg(full_name)) ON CONFLICT DO NOTHING;

-- name: UpdateUser :exec
UPDATE users
SET full_name = sqlc.arg(full_name),
    email = sqlc.arg(email),
    phone = sqlc.arg(phone),
    profile_photo = sqlc.arg(profile_photo),
    status = sqlc.arg(status)
WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email);

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = sqlc.arg(phone);

-- name: GetTeacherById :one
SELECT * FROM users
WHERE id = sqlc.arg(id) AND auth_role = 2;

-- name: HandleUserStatus :exec
UPDATE users
SET status = 0
WHERE id = sqlc.arg(id)
AND status = sqlc.arg(status);

-- name: VerifyLearner :exec
UPDATE verification_learners
SET verified_by = sqlc.arg(verified_by),
    status = sqlc.arg(status),
    verified_at = NOW()
WHERE learner_id = sqlc.arg(learner_id);

-- name: GetLearnerVerificationByUserId :one
SELECT * FROM verification_learners
WHERE learner_id = sqlc.arg(learner_id);

-- name: GetVerifiedLearnersByLearnerId :one
SELECT u.* FROM
users u JOIN verification_learners vls ON u.id = vls.learner_id
WHERE vls.learner_id = sqlc.arg(learner_id) AND vls.status = 1;

-- name: UpdateLearner :exec
INSERT INTO verification_learners (learner_id, school_id, type, image_link, id)
VALUES (sqlc.arg(learner_id), sqlc.arg(school_id), sqlc.arg(type), sqlc.arg(image_link), uuid_generate_v4())
ON CONFLICT (learner_id)
DO
UPDATE SET school_id = sqlc.arg(school_id), type = sqlc.arg(type), image_link = sqlc.arg(image_link), status = 0
WHERE verification_learners.status = 0 OR verification_learners.status = 2;