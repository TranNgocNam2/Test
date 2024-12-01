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
		users.GET("/current", handlers.GetCurrentUser())
		users.GET("", handlers.GetUsers())
		users.GET("/verifications", handlers.GetVerificationUsers())
		users.GET("/:id", handlers.GetUserById())
		users.PUT("/:id", handlers.UpdateUser())
		users.PUT("/verifications/:verificationId", handlers.VerifyUser())
		users.PUT("/:id/handle", handlers.HandleUser())
		users.POST("/learners", handlers.CreateLearner())
		users.PUT("/:id/learners", handlers.UpdateLearner())
	}
}
