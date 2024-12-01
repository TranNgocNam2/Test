package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func ClassRoutes(router *gin.Engine, app *app.Application) {
	classCore := class.NewCore(app)
	handlers := New(classCore)
	classes := router.Group("/classes")
	{
		classes.POST("", handlers.CreateClass())
		classes.POST("/:id/learners/import", handlers.ImportLearners())
		classes.GET("/:id", handlers.GetClassById())
		classes.PUT("/:id", handlers.UpdateClass())
		classes.GET("", handlers.GetClassesByManager())
		classes.GET("/learners", handlers.GetClassesByLearner())
		classes.DELETE("/:id", handlers.DeleteClass())
		classes.PUT("/:id/slots", handlers.UpdateClassSlot())
		classes.POST("/slots/teachers", handlers.CheckTeacherAvailable())
		classes.GET("/teachers", handlers.GetClassesByTeacher())
		classes.PUT("/:id/meeting", handlers.UpdateMeetingLink())
	}
}
