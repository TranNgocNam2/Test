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
	staff, err := queries.GetUserByID(ctx, ctx.GetHeader(header.XUserId))
	if err != nil || (staff.AuthRole != role.MANAGER && staff.AuthRole != role.ADMIN) {
		return "", ErrInvalidUser
	}

	return staff.ID, nil
}
