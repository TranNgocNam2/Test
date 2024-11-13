package middleware

import (
	"Backend/business/db/sqlc"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/enum/role"
)

func AuthorizeStaff(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	staff, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || staff.AuthRole == role.LEARNER || staff.AuthRole == role.TEACHER {
		return "", ErrInvalidUser
	}

	return staff.ID, nil
}

func AuthorizeVerifiedLearner(ctx *gin.Context, queries *sqlc.Queries) (sqlc.User, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return sqlc.User{}, ErrInvalidUser
	}

	learner, err := queries.GetVerifiedLearnerById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil {
		return sqlc.User{}, ErrInvalidUser
	}

	return learner, nil
}

func AuthorizeTeacher(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	teacher, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || teacher.AuthRole != role.TEACHER {
		return "", ErrInvalidUser
	}

	return teacher.ID, nil
}
func AuthorizeWithoutLearner(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	user, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || user.AuthRole == role.LEARNER {
		return "", ErrInvalidUser
	}

	return user.ID, nil
}

func AuthorizeUser(ctx *gin.Context, queries *sqlc.Queries) error {
	if ctx.GetHeader(header.XUserId) == "" {
		return ErrInvalidUser
	}

	_, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil {
		return ErrInvalidUser
	}

	return nil
}
