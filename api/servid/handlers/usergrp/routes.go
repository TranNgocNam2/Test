package usergrp

import (
	"Backend/business/core/user"
	"Backend/business/core/user/userdb"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, app *app.Application) {
	userCore := user.NewCore(userdb.NewStore(app.Db, app.Queries))
	handlers := New(userCore)

	users := router.Group("/users")
	{
		users.POST("", handlers.CreateUser())
		users.GET("/:id", handlers.GetUserByID())
	}
}
