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
INSERT INTO users (id, email, auth_role, full_name)
VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING
`

type CreateUserParams struct {
	ID       string  `db:"id" json:"id"`
	Email    string  `db:"email" json:"email"`
	AuthRole int16   `db:"auth_role" json:"authRole"`
	FullName *string `db:"full_name" json:"fullName"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.AuthRole,
		arg.FullName,
	)
	return err
}

const getLearnerVerificationByUserId = `-- name: GetLearnerVerificationByUserId :one
SELECT id, school_id, learner_id, image_link, status, verified_by, type, verified_at FROM verification_learners
WHERE learner_id = $1
`

func (q *Queries) GetLearnerVerificationByUserId(ctx context.Context, learnerID string) (VerificationLearner, error) {
	row := q.db.QueryRow(ctx, getLearnerVerificationByUserId, learnerID)
	var i VerificationLearner
	err := row.Scan(
		&i.ID,
		&i.SchoolID,
		&i.LearnerID,
		&i.ImageLink,
		&i.Status,
		&i.VerifiedBy,
		&i.Type,
		&i.VerifiedAt,
	)
	return i, err
}

const getTeacherById = `-- name: GetTeacherById :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status FROM users
WHERE id = $1 AND auth_role = 2
`

func (q *Queries) GetTeacherById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getTeacherById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status FROM users
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
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status FROM users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id string) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status FROM users
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
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
	)
	return i, err
}

const getVerifiedLearnersByLearnerId = `-- name: GetVerifiedLearnersByLearnerId :one
SELECT u.id, u.full_name, u.email, u.phone, u.auth_role, u.profile_photo, u.status FROM
users u JOIN verification_learners vls ON u.id = vls.learner_id
WHERE vls.learner_id = $1 AND vls.status = 1
`

func (q *Queries) GetVerifiedLearnersByLearnerId(ctx context.Context, learnerID string) (User, error) {
	row := q.db.QueryRow(ctx, getVerifiedLearnersByLearnerId, learnerID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
	)
	return i, err
}

const handleUserStatus = `-- name: HandleUserStatus :exec
UPDATE users
SET status = 0
WHERE id = $1
AND status = $2
`

type HandleUserStatusParams struct {
	ID     string `db:"id" json:"id"`
	Status int32  `db:"status" json:"status"`
}

func (q *Queries) HandleUserStatus(ctx context.Context, arg HandleUserStatusParams) error {
	_, err := q.db.Exec(ctx, handleUserStatus, arg.ID, arg.Status)
	return err
}

const updateLearner = `-- name: UpdateLearner :exec
INSERT INTO verification_learners (learner_id, school_id, type, image_link, id)
VALUES ($1, $2, $3, $4, uuid_generate_v4())
ON CONFLICT (learner_id)
DO
UPDATE SET school_id = $2, type = $3, image_link = $4, status = 0
WHERE status = 0 OR status = 2
`

type UpdateLearnerParams struct {
	LearnerID string    `db:"learner_id" json:"learnerId"`
	SchoolID  uuid.UUID `db:"school_id" json:"schoolId"`
	Type      int16     `db:"type" json:"type"`
	ImageLink []string  `db:"image_link" json:"imageLink"`
}

func (q *Queries) UpdateLearner(ctx context.Context, arg UpdateLearnerParams) error {
	_, err := q.db.Exec(ctx, updateLearner,
		arg.LearnerID,
		arg.SchoolID,
		arg.Type,
		arg.ImageLink,
	)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET full_name = $1,
    email = $2,
    phone = $3,
    profile_photo = $4,
    status = $5
WHERE id = $6
`

type UpdateUserParams struct {
	FullName     *string `db:"full_name" json:"fullName"`
	Email        string  `db:"email" json:"email"`
	Phone        *string `db:"phone" json:"phone"`
	ProfilePhoto *string `db:"profile_photo" json:"profilePhoto"`
	Status       int32   `db:"status" json:"status"`
	ID           string  `db:"id" json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.FullName,
		arg.Email,
		arg.Phone,
		arg.ProfilePhoto,
		arg.Status,
		arg.ID,
	)
	return err
}

const verifyLearner = `-- name: VerifyLearner :exec
UPDATE verification_learners
SET verified_by = $1,
    status = $2,
    verified_at = NOW()
WHERE learner_id = $3
`

type VerifyLearnerParams struct {
	VerifiedBy *string `db:"verified_by" json:"verifiedBy"`
	Status     int16   `db:"status" json:"status"`
	LearnerID  string  `db:"learner_id" json:"learnerId"`
}

func (q *Queries) VerifyLearner(ctx context.Context, arg VerifyLearnerParams) error {
	_, err := q.db.Exec(ctx, verifyLearner, arg.VerifiedBy, arg.Status, arg.LearnerID)
	return err
}
