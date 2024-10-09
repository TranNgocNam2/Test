// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, email, auth_role)
VALUES ($1, $2, $3) ON CONFLICT DO NOTHING
`

type CreateUserParams struct {
	ID       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	AuthRole int16  `db:"auth_role" json:"authRole"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser, arg.ID, arg.Email, arg.AuthRole)
	return err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, email, phone, gender, auth_role, profile_photo, status, school_id FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Gender,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
		&i.SchoolID,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, full_name, email, phone, gender, auth_role, profile_photo, status, school_id FROM users
WHERE id = $1
`

func (q *Queries) GetUserByID(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Gender,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
		&i.SchoolID,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, full_name, email, phone, gender, auth_role, profile_photo, status, school_id FROM users
WHERE phone = $1
`

func (q *Queries) GetUserByPhone(ctx context.Context, phone *string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByPhone, phone)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Gender,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
		&i.SchoolID,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET full_name = $1,
    email = $2,
    phone = $3,
    gender = $4,
    school_id = $5,
    profile_photo = $6
WHERE id = $7
`

type UpdateUserParams struct {
	FullName     *string       `db:"full_name" json:"fullName"`
	Email        string        `db:"email" json:"email"`
	Phone        *string       `db:"phone" json:"phone"`
	Gender       *int16        `db:"gender" json:"gender"`
	SchoolID     uuid.NullUUID `db:"school_id" json:"schoolId"`
	ProfilePhoto *string       `db:"profile_photo" json:"profilePhoto"`
	ID           string        `db:"id" json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.FullName,
		arg.Email,
		arg.Phone,
		arg.Gender,
		arg.SchoolID,
		arg.ProfilePhoto,
		arg.ID,
	)
	return err
}
