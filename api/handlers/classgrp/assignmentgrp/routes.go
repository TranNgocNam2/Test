package assignmentgrp

import (
	"Backend/business/core/class/assignment"
	"Backend/internal/app"

	"github.com/gin-gonic/gin"
)

func AssignmentRoutes(router *gin.Engine, app *app.Application) {
	assignment := assignment.NewCore(app)
	handlers := New(assignment)

	classes := router.Group("/classes")
	{
		classes.POST("/:id/assignments", handlers.CreateAssignment())
		classes.PUT("/:id/assignments/:assignmentId", handlers.UpdateAssignment())
		classes.GET("/:id/assignments", handlers.GetAssignments())
	}

	assignments := router.Group("/assignments")
	{
		assignments.DELETE("/:id", handlers.DeleteAssignment())
		assignments.GET("/:id", handlers.GetById())
		assignments.PUT("/:id/learners/:learnerId", handlers.GradeAssignment())
		assignments.POST("/:id/learnerAssignments", handlers.SubmitAssignment())
	}

	learners := router.Group("/learnerAssignments")
	{
		learners.GET("/:id", handlers.GetLearnerAssignment())
	}
}
