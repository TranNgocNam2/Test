package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func ProgramRoutes(router *gin.Engine, app *app.Application) {
	programCore := program.NewCore(app)
	handlers := New(programCore)
	programs := router.Group("/programs")
	{
		programs.POST("", handlers.CreateProgram())
		programs.GET("", handlers.GetPrograms())
		programs.PUT("/:id", handlers.UpdateProgram())
		programs.DELETE("/:id", handlers.DeleteProgram())
	}
}
