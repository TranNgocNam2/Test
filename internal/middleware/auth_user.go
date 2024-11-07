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
