package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, app *app.Application) {
	userCore := user.NewCore(app)
	handlers := New(userCore)

	users := router.Group("/users")
	{
		users.POST("", handlers.CreateUser())
		users.GET("/:id", handlers.GetUserByID())
		users.PUT("/:id", handlers.UpdateUser())
		users.PUT("/:id/verify", handlers.VerifyUser())
	}
}
