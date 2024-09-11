package usergrp

import (
	"Backend/internal/platform/app"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, app *app.Application) {
	router.GET("/users", GetAllUsersHandler(app))
	router.POST("/users", CreateUserHandler(app))
}
