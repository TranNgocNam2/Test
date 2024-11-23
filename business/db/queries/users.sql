-- name: GetUserById :one
SELECT * FROM users
WHERE id = sqlc.arg(id);

-- name: CreateUser :exec
INSERT INTO users (id, email, auth_role, full_name, is_verified)
VALUES (sqlc.arg(id), sqlc.arg(email), sqlc.arg(auth_role), sqlc.arg(full_name),
        sqlc.arg(is_verified)) ON CONFLICT DO NOTHING;

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
SET status = sqlc.arg(status)
WHERE id = sqlc.arg(id);

-- name: VerifyLearner :exec
UPDATE verification_learners
SET verified_by = sqlc.arg(verified_by),
    status = sqlc.arg(status),
    note = sqlc.arg(note)::text,
    verified_at = NOW()
WHERE learner_id = sqlc.arg(learner_id)
AND id = sqlc.arg(id);

-- name: UpdateVerification :exec
UPDATE users
SET is_verified = sqlc.arg(is_verified),
    school_id = sqlc.arg(school_id),
    type = sqlc.arg(type)
WHERE id = sqlc.arg(id);

-- name: GetLearnerVerificationById :one
SELECT *
FROM verification_learners
WHERE id = sqlc.arg(id);

-- name: GetLearnerVerificationByLearnerId :one
SELECT *
FROM verification_learners
WHERE learner_id = sqlc.arg(learner_id)
  AND status = sqlc.arg(status);


-- name: GetVerifiedLearnersByLearnerId :one
SELECT * FROM users
WHERE id = sqlc.arg(id) AND is_verified = true AND status = sqlc.arg(status);

-- name: CreateVerificationRequest :one
INSERT INTO verification_learners (learner_id, school_id, type, image_link, id, status)
VALUES (sqlc.arg(learner_id), sqlc.arg(school_id), sqlc.arg(type),
        sqlc.arg(image_link), sqlc.arg(id), sqlc.arg(status)) RETURNING id;

-- name: GetVerificationLearners :many
SELECT u.id AS user_id, u.full_name, u.email,
       vl.id, vl.image_link::text AS image_link, vl.type, vl.status, vl.note,
       s.id AS school_id, s.name AS school_name
FROM users u
JOIN verification_learners vl ON u.id = vl.learner_id
JOIN schools s ON vl.school_id = s.id
WHERE vl.learner_id = sqlc.arg(learner_id);