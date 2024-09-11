package main

import (
	"Backend/api/servid/handlers/usergrp"
	"Backend/api/servid/routes"
	"Backend/internal/platform/app"
	"Backend/internal/platform/config"
	"Backend/internal/platform/db"
	"Backend/internal/platform/db/ent"
	"Backend/internal/platform/log"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/file"
	"gitlab.com/innovia69420/kit/logger"
	"io"
	"os"
)

var workingDirectory string

func init() {
	var err error
	workingDirectory, err = file.WorkingDirectory()
	if err != nil {
		fmt.Println("Error getting root path:", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
}

func main() {
	router := gin.Default()
	//router.Use(middleware.RecoverPanic())

	cfg, _ := config.LoadAllAppConfig(workingDirectory)

	//Config Cors
	//corsConfig := cors.Config{
	//	AllowOrigins:     []string{cfg.CorsOrigin},
	//	AllowMethods:     enum.CorsAllowMethods,
	//	AllowHeaders:     enum.CorsAllowHeaders,
	//	AllowCredentials: true,
	//	MaxAge:           enum.CorsMaxAge,
	//}
	//router.Use(cors.New(corsConfig))
	//Set up log
	zapLog := logger.Get(workingDirectory)
	router.Use(log.RequestLogger(zapLog))

	ctx := context.Background()

	client := db.ConnectDB(ctx, cfg.DatabaseUrl, zapLog)
	//Create Schema
	app := app.Application{
		Config:    cfg,
		EntClient: client,
		Logger:    zapLog,
	}

	if err := client.Schema.Create(ctx); err != nil {
		logger.StartUpError(app.Logger, message.FailedCreateEntSchema)
	}

	//Connect DB and close connection

	defer func(client *ent.Client) {
		_ = client.Close()
	}(client)

	// Load all routes
	LoadRoutes(router, &app)

	serverAddr := fmt.Sprintf("%s:%d", app.Config.Host, app.Config.Port)
	app.Logger.Info("Server is listening on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logger.StartUpError(app.Logger, message.FailedStartApplication)
	}

}

func LoadRoutes(router *gin.Engine, app *app.Application) {
	routes.ExampleRoutes(router)
	usergrp.UserRoutes(router, app)
}
