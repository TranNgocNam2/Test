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
		learners.POST("/specializations/:specializationId", handlers.AddLearnerToSpecialization())
		learners.PUT("/classes/:classId/attendance", handlers.SubmitAttendance())
	}
}
