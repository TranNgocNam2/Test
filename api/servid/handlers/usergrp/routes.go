package usergrp

import (
	"Backend/business/core/user"
	"Backend/business/core/user/userdb"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, app *app.Application) {
	userCore := user.NewCore(userdb.NewStore(app.Db, app.Queries))
	handlers := New(userCore)

	router.POST("/users", handlers.CreateUser())
	//router.PUT("/schools/:id", handlers.UpdateSchool())
	//router.DELETE("/schools/:id", handlers.DeleteSchool())
	router.GET("/users/:id", handlers.GetUserByID())
	//router.GET("/districts/:id/schools", handlers.GetSchoolsByDistrict())
	//router.GET("/provinces", handlers.GetProvinces())
	//router.GET("/provinces/:id/districts", handlers.GetDistrictsByProvince())
}
