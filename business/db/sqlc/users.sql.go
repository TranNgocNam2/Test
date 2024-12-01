// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createLearner = `-- name: CreateLearner :exec
INSERT INTO users (id, email, auth_role, full_name, is_verified, school_id)
VALUES ($1, $2, $3, $4,
        $5, $6) ON CONFLICT DO NOTHING
`

type CreateLearnerParams struct {
	ID         string     `db:"id" json:"id"`
	Email      string     `db:"email" json:"email"`
	AuthRole   int16      `db:"auth_role" json:"authRole"`
	FullName   *string    `db:"full_name" json:"fullName"`
	IsVerified bool       `db:"is_verified" json:"isVerified"`
	SchoolID   *uuid.UUID `db:"school_id" json:"schoolId"`
}

func (q *Queries) CreateLearner(ctx context.Context, arg CreateLearnerParams) error {
	_, err := q.db.Exec(ctx, createLearner,
		arg.ID,
		arg.Email,
		arg.AuthRole,
		arg.FullName,
		arg.IsVerified,
		arg.SchoolID,
	)
	return err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, email, auth_role, full_name, is_verified)
VALUES ($1, $2, $3, $4,
        $5) ON CONFLICT DO NOTHING
`

type CreateUserParams struct {
	ID         string  `db:"id" json:"id"`
	Email      string  `db:"email" json:"email"`
	AuthRole   int16   `db:"auth_role" json:"authRole"`
	FullName   *string `db:"full_name" json:"fullName"`
	IsVerified bool    `db:"is_verified" json:"isVerified"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.Exec(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.AuthRole,
		arg.FullName,
		arg.IsVerified,
	)
	return err
}

const createVerificationRequest = `-- name: CreateVerificationRequest :one
INSERT INTO verification_learners (learner_id, school_id, type, image_link, id, status)
VALUES ($1, $2, $3,
        $4, $5, $6) RETURNING id
`

type CreateVerificationRequestParams struct {
	LearnerID string    `db:"learner_id" json:"learnerId"`
	SchoolID  uuid.UUID `db:"school_id" json:"schoolId"`
	Type      int16     `db:"type" json:"type"`
	ImageLink []string  `db:"image_link" json:"imageLink"`
	ID        uuid.UUID `db:"id" json:"id"`
	Status    int16     `db:"status" json:"status"`
}

func (q *Queries) CreateVerificationRequest(ctx context.Context, arg CreateVerificationRequestParams) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, createVerificationRequest,
		arg.LearnerID,
		arg.SchoolID,
		arg.Type,
		arg.ImageLink,
		arg.ID,
		arg.Status,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getEmailsExcept = `-- name: GetEmailsExcept :one
SELECT STRING_AGG(email, ', ') AS emails
FROM (
         SELECT DISTINCT email
         FROM UNNEST($1::text[]) AS unnested_emails(email)
         EXCEPT
         SELECT email
         FROM users
         WHERE email = ANY($1::text[])
           AND status = $2
           AND is_verified = $3
           AND auth_role = $4
     ) missing_emails
`

type GetEmailsExceptParams struct {
	Emails     []string `db:"emails" json:"emails"`
	Status     int32    `db:"status" json:"status"`
	IsVerified bool     `db:"is_verified" json:"isVerified"`
	AuthRole   int16    `db:"auth_role" json:"authRole"`
}

func (q *Queries) GetEmailsExcept(ctx context.Context, arg GetEmailsExceptParams) ([]byte, error) {
	row := q.db.QueryRow(ctx, getEmailsExcept,
		arg.Emails,
		arg.Status,
		arg.IsVerified,
		arg.AuthRole,
	)
	var emails []byte
	err := row.Scan(&emails)
	return emails, err
}

const getLearnerVerificationById = `-- name: GetLearnerVerificationById :one
SELECT id, school_id, learner_id, image_link, status, verified_by, type, verified_at, note, created_at
FROM verification_learners
WHERE id = $1
`

func (q *Queries) GetLearnerVerificationById(ctx context.Context, id uuid.UUID) (VerificationLearner, error) {
	row := q.db.QueryRow(ctx, getLearnerVerificationById, id)
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
		&i.Note,
		&i.CreatedAt,
	)
	return i, err
}

const getLearnerVerificationByLearnerId = `-- name: GetLearnerVerificationByLearnerId :one
SELECT id, school_id, learner_id, image_link, status, verified_by, type, verified_at, note, created_at
FROM verification_learners
WHERE learner_id = $1
  AND status = $2
`

type GetLearnerVerificationByLearnerIdParams struct {
	LearnerID string `db:"learner_id" json:"learnerId"`
	Status    int16  `db:"status" json:"status"`
}

func (q *Queries) GetLearnerVerificationByLearnerId(ctx context.Context, arg GetLearnerVerificationByLearnerIdParams) (VerificationLearner, error) {
	row := q.db.QueryRow(ctx, getLearnerVerificationByLearnerId, arg.LearnerID, arg.Status)
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
		&i.Note,
		&i.CreatedAt,
	)
	return i, err
}

const getTeacherById = `-- name: GetTeacherById :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status, is_verified, school_id, type FROM users
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
		&i.IsVerified,
		&i.SchoolID,
		&i.Type,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status, is_verified, school_id, type FROM users
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
		&i.IsVerified,
		&i.SchoolID,
		&i.Type,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status, is_verified, school_id, type FROM users
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
		&i.IsVerified,
		&i.SchoolID,
		&i.Type,
	)
	return i, err
}

const getUserByPhone = `-- name: GetUserByPhone :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status, is_verified, school_id, type FROM users
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
		&i.IsVerified,
		&i.SchoolID,
		&i.Type,
	)
	return i, err
}

const getUsersByEmails = `-- name: GetUsersByEmails :many
SELECT id AS ids
FROM users
WHERE email = ANY($1::text[])
  AND status = $2
  AND is_verified = $3
  AND auth_role = $4
`

type GetUsersByEmailsParams struct {
	Emails     []string `db:"emails" json:"emails"`
	Status     int32    `db:"status" json:"status"`
	IsVerified bool     `db:"is_verified" json:"isVerified"`
	AuthRole   int16    `db:"auth_role" json:"authRole"`
}

func (q *Queries) GetUsersByEmails(ctx context.Context, arg GetUsersByEmailsParams) ([]string, error) {
	rows, err := q.db.Query(ctx, getUsersByEmails,
		arg.Emails,
		arg.Status,
		arg.IsVerified,
		arg.AuthRole,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var ids string
		if err := rows.Scan(&ids); err != nil {
			return nil, err
		}
		items = append(items, ids)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVerificationLearners = `-- name: GetVerificationLearners :many
SELECT u.id AS user_id, u.full_name, u.email,
       vl.id, vl.image_link::text AS image_link, vl.type, vl.status, vl.note, vl.created_at,
       s.id AS school_id, s.name AS school_name
FROM users u
JOIN verification_learners vl ON u.id = vl.learner_id
JOIN schools s ON vl.school_id = s.id
WHERE vl.learner_id = $1
`

type GetVerificationLearnersRow struct {
	UserID     string    `db:"user_id" json:"userId"`
	FullName   *string   `db:"full_name" json:"fullName"`
	Email      string    `db:"email" json:"email"`
	ID         uuid.UUID `db:"id" json:"id"`
	ImageLink  string    `db:"image_link" json:"imageLink"`
	Type       int16     `db:"type" json:"type"`
	Status     int16     `db:"status" json:"status"`
	Note       *string   `db:"note" json:"note"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
	SchoolID   uuid.UUID `db:"school_id" json:"schoolId"`
	SchoolName string    `db:"school_name" json:"schoolName"`
}

func (q *Queries) GetVerificationLearners(ctx context.Context, learnerID string) ([]GetVerificationLearnersRow, error) {
	rows, err := q.db.Query(ctx, getVerificationLearners, learnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetVerificationLearnersRow
	for rows.Next() {
		var i GetVerificationLearnersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FullName,
			&i.Email,
			&i.ID,
			&i.ImageLink,
			&i.Type,
			&i.Status,
			&i.Note,
			&i.CreatedAt,
			&i.SchoolID,
			&i.SchoolName,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVerifiedLearnersByLearnerId = `-- name: GetVerifiedLearnersByLearnerId :one
SELECT id, full_name, email, phone, auth_role, profile_photo, status, is_verified, school_id, type FROM users
WHERE id = $1 AND is_verified = true AND status = $2
`

type GetVerifiedLearnersByLearnerIdParams struct {
	ID     string `db:"id" json:"id"`
	Status int32  `db:"status" json:"status"`
}

func (q *Queries) GetVerifiedLearnersByLearnerId(ctx context.Context, arg GetVerifiedLearnersByLearnerIdParams) (User, error) {
	row := q.db.QueryRow(ctx, getVerifiedLearnersByLearnerId, arg.ID, arg.Status)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.AuthRole,
		&i.ProfilePhoto,
		&i.Status,
		&i.IsVerified,
		&i.SchoolID,
		&i.Type,
	)
	return i, err
}

const handleUserStatus = `-- name: HandleUserStatus :exec
UPDATE users
SET status = $1
WHERE id = $2
`

type HandleUserStatusParams struct {
	Status int32  `db:"status" json:"status"`
	ID     string `db:"id" json:"id"`
}

func (q *Queries) HandleUserStatus(ctx context.Context, arg HandleUserStatusParams) error {
	_, err := q.db.Exec(ctx, handleUserStatus, arg.Status, arg.ID)
	return err
}

const updateLearner = `-- name: UpdateLearner :exec
UPDATE users
SET school_id = $1,
    type = $2
WHERE id = $3
`

type UpdateLearnerParams struct {
	SchoolID *uuid.UUID `db:"school_id" json:"schoolId"`
	Type     *int16     `db:"type" json:"type"`
	ID       string     `db:"id" json:"id"`
}

func (q *Queries) UpdateLearner(ctx context.Context, arg UpdateLearnerParams) error {
	_, err := q.db.Exec(ctx, updateLearner, arg.SchoolID, arg.Type, arg.ID)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET full_name = $1,
    profile_photo = $2
WHERE id = $3
`

type UpdateUserParams struct {
	FullName     *string `db:"full_name" json:"fullName"`
	ProfilePhoto *string `db:"profile_photo" json:"profilePhoto"`
	ID           string  `db:"id" json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser, arg.FullName, arg.ProfilePhoto, arg.ID)
	return err
}

const updateVerification = `-- name: UpdateVerification :exec
UPDATE users
SET is_verified = $1,
    school_id = $2,
    type = $3
WHERE id = $4
`

type UpdateVerificationParams struct {
	IsVerified bool       `db:"is_verified" json:"isVerified"`
	SchoolID   *uuid.UUID `db:"school_id" json:"schoolId"`
	Type       *int16     `db:"type" json:"type"`
	ID         string     `db:"id" json:"id"`
}

func (q *Queries) UpdateVerification(ctx context.Context, arg UpdateVerificationParams) error {
	_, err := q.db.Exec(ctx, updateVerification,
		arg.IsVerified,
		arg.SchoolID,
		arg.Type,
		arg.ID,
	)
	return err
}

const verifyLearner = `-- name: VerifyLearner :exec
UPDATE verification_learners
SET verified_by = $1,
    status = $2,
    note = $3::text,
    verified_at = NOW()
WHERE learner_id = $4
AND id = $5
`

type VerifyLearnerParams struct {
	VerifiedBy *string   `db:"verified_by" json:"verifiedBy"`
	Status     int16     `db:"status" json:"status"`
	Note       string    `db:"note" json:"note"`
	LearnerID  string    `db:"learner_id" json:"learnerId"`
	ID         uuid.UUID `db:"id" json:"id"`
}

func (q *Queries) VerifyLearner(ctx context.Context, arg VerifyLearnerParams) error {
	_, err := q.db.Exec(ctx, verifyLearner,
		arg.VerifiedBy,
		arg.Status,
		arg.Note,
		arg.LearnerID,
		arg.ID,
	)
	return err
}
