package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SchoolRoutes(router *gin.Engine, app *app.Application) {
	schoolCore := school.NewCore(app)
	hdl := New(schoolCore)

	router.GET("/provinces", hdl.GetProvinces())
	router.GET("/provinces/:province_id/districts", hdl.GetDistrictsByProvince())
}
