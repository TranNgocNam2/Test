package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SpecializationRoutes(router *gin.Engine, app *app.Application) {
	specCore := specialization.NewCore(app)
	handlers := New(specCore)
	specializations := router.Group("/specializations")
	{
		specializations.POST("", handlers.CreateSpecialization())
		specializations.GET("/:id", handlers.GetSpecializationById())
		specializations.GET("", handlers.GetSpecializations())
		specializations.PUT("/:id", handlers.UpdateSpecialization())
		specializations.DELETE("/:id", handlers.DeleteSpecialization())
	}
}
