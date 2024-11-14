package teachergrp

import (
	"Backend/business/core/teacher"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func TeacherRoutes(router *gin.Engine, app *app.Application) {
	teacherCore := teacher.NewCore(app)
	handlers := New(teacherCore)
	teachers := router.Group("/teachers")
	{
		teachers.PUT("/slots/:slotId", handlers.GenerateAttendanceCode())
	}
}
