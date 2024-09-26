-- name: GetUserByID :one
SELECT * FROM users
WHERE id = sqlc.arg(id);

-- name: CreateUser :exec
INSERT INTO users (id, full_name, email, phone, gender, profile_photo, school_id, role)
VALUES (sqlc.arg(id), sqlc.arg(full_name), sqlc.arg(email), sqlc.arg(phone),
        sqlc.arg(gender), sqlc.arg(profile_photo), sqlc.arg(school_id), sqlc.arg(role));

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = sqlc.arg(email);

-- name: GetUserByPhone :one
SELECT * FROM users
WHERE phone = sqlc.arg(phone);