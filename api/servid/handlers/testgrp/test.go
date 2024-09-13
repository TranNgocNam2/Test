package testgrp

import (
	"Backend/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/web"
)

func GetAccountsHandler(app *app.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accounts, err := app.Queries.GetAccounts(ctx)
		if err != nil {
			web.BadRequestError(ctx, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, accounts)
	}
}
