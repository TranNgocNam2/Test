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
		classes.GET("/:id", handlers.GetClassById())
		classes.PUT("/:id", handlers.UpdateClass())
		classes.GET("", handlers.GetClassesByManager())
		classes.DELETE("/:id", handlers.DeleteClass())
		classes.PUT("/:id/teachers", handlers.UpdateClassTeacher())
		classes.PUT("/:id/slots", handlers.UpdateClassSlot())
		classes.GET("/slots/teachers", handlers.CheckTeacherConflict())
	}
}
