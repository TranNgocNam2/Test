package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/business/core/school/schooldb"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SchoolRoutes(router *gin.Engine, app *app.Application) {
	schoolCore := school.NewCore(schooldb.NewStore(app.Db, app.Queries, app.Logger))
	handlers := New(schoolCore)

	router.POST("/schools", handlers.CreateSchool())
	router.PUT("/schools/:id", handlers.UpdateSchool())
	router.DELETE("/schools/:id", handlers.DeleteSchool())
	router.GET("/schools/:id", handlers.GetSchoolByID())
	router.GET("/schools", handlers.GetSchoolPaginated())
	router.GET("/districts/:id/schools", handlers.GetSchoolsByDistrict())
	router.GET("/provinces", handlers.GetProvinces())
	router.GET("/provinces/:id/districts", handlers.GetDistrictsByProvince())
}
