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
		classes.PUT("/:id", handlers.UpdateClass())
		classes.DELETE("/:id", handlers.DeleteClass())
	}
}
