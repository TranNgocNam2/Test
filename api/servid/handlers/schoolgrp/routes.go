package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/business/core/school/schooldb"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SchoolRoutes(router *gin.Engine, app *app.Application) {
	schoolCore := school.NewCore(schooldb.NewStore(app.Db, app.Queries))
	handlers := New(schoolCore)

	schools := router.Group("/schools")
	{
		schools.GET("", handlers.CreateSchool())
		schools.GET("/:id", handlers.GetSchoolByID())
		schools.DELETE("/:id", handlers.DeleteSchool())
		schools.PUT("/:id", handlers.UpdateSchool())
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
