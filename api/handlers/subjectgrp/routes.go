package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/app"

	"github.com/gin-gonic/gin"
)

func SubjectRoutes(router *gin.Engine, app *app.Application) {
	core := subject.NewCore(app)
	handlers := New(core)
	subjects := router.Group("/subjects")
	{
		subjects.POST("", handlers.CreateSubject())
		subjects.PUT("/:id", handlers.UpdateSubject())
		subjects.GET("/:id", handlers.GetSubjectById())
		subjects.GET("", handlers.GetSubjects())
		subjects.DELETE("/:id", handlers.DeleteSubject())
	}
}
