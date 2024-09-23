package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SchoolRoutes(router *gin.Engine, app *app.Application) {
	schoolCore := school.NewCore(app)
	handlers := New(schoolCore)

	router.POST("/schools", handlers.CreateSchool())
	router.DELETE("/schools/:id", handlers.DeleteSchool())
	router.GET("/provinces", handlers.GetProvinces())
	router.GET("/provinces/:id/districts", handlers.GetDistrictsByProvince())
}
