package middleware

import (
	"Backend/business/db/sqlc"
	"Backend/internal/common/status"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/enum/role"
)

func AuthorizeStaff(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	staff, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || staff.AuthRole == role.LEARNER || staff.AuthRole == role.TEACHER || status.User(staff.Status) != status.Valid {
		return "", ErrInvalidUser
	}

	return staff.ID, nil
}

func AuthorizeVerifiedLearner(ctx *gin.Context, queries *sqlc.Queries) (*sqlc.User, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return nil, ErrInvalidUser
	}

	learner, err := queries.GetVerifiedLearnersByLearnerId(ctx,
		sqlc.GetVerifiedLearnersByLearnerIdParams{
			ID:     ctx.GetHeader(header.XUserId),
			Status: int32(status.Valid),
		})
	if err != nil {
		return nil, ErrInvalidUser
	}

	return &learner, nil
}

func AuthorizeTeacher(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	teacher, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	fmt.Println(teacher.Status)
	if err != nil || teacher.AuthRole != role.TEACHER || status.User(teacher.Status) != status.Valid {
		fmt.Println(err.Error())
		return "", ErrInvalidUser
	}

	return teacher.ID, nil
}
func AuthorizeWithoutLearner(ctx *gin.Context, queries *sqlc.Queries) (string, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return "", ErrInvalidUser
	}

	user, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || user.AuthRole == role.LEARNER || status.User(user.Status) != status.Valid {
		return "", ErrInvalidUser
	}

	return user.ID, nil
}

func AuthorizeUser(ctx *gin.Context, queries *sqlc.Queries) (*sqlc.User, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return nil, ErrInvalidUser
	}

	user, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || status.User(user.Status) != status.Valid {
		return nil, ErrInvalidUser
	}

	return &user, nil
}

func AuthorizeAdmin(ctx *gin.Context, queries *sqlc.Queries) (*sqlc.User, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return nil, ErrInvalidUser
	}

	admin, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || admin.AuthRole != role.ADMIN {
		return nil, ErrInvalidUser
	}

	return &admin, nil
}

func AuthorizeLearner(ctx *gin.Context, queries *sqlc.Queries) (*sqlc.User, error) {
	if ctx.GetHeader(header.XUserId) == "" {
		return nil, ErrInvalidUser
	}

	learner, err := queries.GetUserById(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || learner.AuthRole != role.LEARNER || status.User(learner.Status) != status.Valid {
		return nil, ErrInvalidUser
	}

	return &learner, nil
}
