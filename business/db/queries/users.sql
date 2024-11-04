-- name: GetUserById :one
SELECT * FROM users
WHERE id = sqlc.arg(id);

-- name: CreateUser :exec
INSERT INTO users (id, email, auth_role)
VALUES (sqlc.arg(id), sqlc.arg(email), sqlc.arg(auth_role)) ON CONFLICT DO NOTHING;

-- name: UpdateUser :exec
UPDATE users
SET full_name = sqlc.arg(full_name),
    email = sqlc.arg(email),
    phone = sqlc.arg(phone),
    gender = sqlc.arg(gender),
    school_id = sqlc.arg(school_id),
    profile_photo = sqlc.arg(profile_photo)
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