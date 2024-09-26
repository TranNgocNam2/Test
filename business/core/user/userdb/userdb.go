package userdb

import (
	"Backend/business/core/user"
	"Backend/business/db/sqlc"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"net/mail"
)

var (
	ErrNotFound     = errors.New("user not found")
	FailedToGetUser = errors.New("failed to get user")
	FailedToCreate  = errors.New("failed to create user")
)

type Store struct {
	db      *sqlx.DB
	queries *sqlc.Queries
}

func NewStore(db *sqlx.DB, queries *sqlc.Queries) *Store {
	return &Store{
		db:      db,
		queries: queries,
	}
}

func (s *Store) Create(ctx *gin.Context, user user.User) error {
	newUserDB := sqlc.CreateUserParams{
		ID:           user.ID,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        user.Phone,
		Gender:       user.Gender,
		ProfilePhoto: user.Photo,
		SchoolID: uuid.NullUUID{
			UUID:  user.School.ID,
			Valid: user.School.ID != uuid.Nil,
		},
		Role: user.Role,
	}
	if err := s.queries.CreateUser(ctx, newUserDB); err != nil {
		return FailedToCreate
	}

	return nil
}

func (s *Store) GetByEmail(ctx *gin.Context, email mail.Address) (user.User, error) {
	dbUser, err := s.queries.GetUserByEmail(ctx, email.Address)
	if err != nil {
		return user.User{}, ErrNotFound
	}
	return toCoreUser(dbUser), nil
}

func (s *Store) GetByPhone(ctx *gin.Context, phone string) (user.User, error) {
	dbUser, err := s.queries.GetUserByPhone(ctx, phone)
	if err != nil {
		return user.User{}, ErrNotFound
	}
	return toCoreUser(dbUser), nil
}

func (s *Store) GetByID(ctx *gin.Context, id string) (user.User, error) {
	dbUser, err := s.queries.GetUserByID(ctx, id)
	if err != nil {
		return user.User{}, ErrNotFound
	}

	if !dbUser.SchoolID.Valid {
		return toCoreUser(dbUser), nil
	}

	dbSchool, err := s.queries.GetSchoolByID(ctx, dbUser.SchoolID.UUID)
	if err != nil {
		return user.User{}, FailedToGetUser
	}

	coreUser := toCoreUser(dbUser)
	coreUser.School.ID = dbSchool.ID
	coreUser.School.Name = dbSchool.Name

	return coreUser, nil
}
