-- name: GetUserByID :one
SELECT * FROM users
WHERE id = sqlc.arg(id);

-- name: GetLearnerByID :one
SELECT * FROM learners
WHERE id = sqlc.arg(id);

-- name: CreateUser :exec
INSERT INTO users (id, full_name, email, phone, gender, profile_photo, auth_role)
VALUES (sqlc.arg(id), sqlc.arg(full_name), sqlc.arg(email), sqlc.arg(phone),
        sqlc.arg(gender), sqlc.arg(profile_photo), sqlc.arg(auth_role));

-- name: CreateLeaner :exec
INSERT INTO learners (id, school_id)
VALUES (sqlc.arg(id), sqlc.arg(school_id));

-- name: CreateStaff :exec
INSERT INTO staffs (id, role, created_by)
VALUES (sqlc.arg(id), sqlc.arg(role), sqlc.arg(created_by));

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email);

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = sqlc.arg(phone);