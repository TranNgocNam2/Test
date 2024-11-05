package learnergrp

import (
	"Backend/business/core/learner"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func LearnerRoutes(router *gin.Engine, app *app.Application) {
	learnerCore := learner.NewCore(app)
	handlers := New(learnerCore)
	learners := router.Group("/learners")
	{
		learners.POST("", handlers.AddLearnerToClass())
		//learners.GET("/:id", handlers.GetSpecializationById())
		//learners.GET("", handlers.GetSpecializations())
		//learners.PUT("/:id", handlers.UpdateSpecialization())
		//learners.DELETE("/:id", handlers.DeleteSpecialization())
	}
}
