package testgrp

import (
	"Backend/internal/app"

	"github.com/gin-gonic/gin"
)

func AccountRoutes(router *gin.Engine, app *app.Application) {
	router.GET("/accounts", GetAccountsHandler(app))
}
