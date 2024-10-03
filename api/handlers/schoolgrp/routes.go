package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SchoolRoutes(router *gin.Engine, app *app.Application) {
	schoolCore := school.NewCore(app)
	handlers := New(schoolCore)

	schools := router.Group("/schools")
	{
		schools.POST("", handlers.CreateSchool())
		schools.GET("/:id", handlers.GetSchoolByID())
		schools.DELETE("/:id", handlers.DeleteSchool())
		schools.PUT("/:id", handlers.UpdateSchool())
		schools.GET("", handlers.GetSchoolPaginated())
	}

	provinces := router.Group("/provinces")
	{
		provinces.GET("", handlers.GetProvinces())
		provinces.GET("/:id/districts", handlers.GetDistrictsByProvince())
	}

	districts := router.Group("/districts")
	{
		districts.GET("/:id/schools", handlers.GetSchoolsByDistrict())

	}
}
