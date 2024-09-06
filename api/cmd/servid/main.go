package main

import (
	"Backend/api/cmd/servid/routes"
	"Backend/api/internal/platform/app"
	"Backend/api/internal/platform/config"
	"Backend/api/internal/platform/db"
	"Backend/api/internal/platform/db/ent"
	"Backend/api/internal/platform/logger"
	"Backend/kit/enum"
	"Backend/kit/file"
	"Backend/kit/log"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

var workingDirectory string

func init() {
	var err error
	workingDirectory, err = file.GetWorkingDirectory()
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
	fmt.Println(workingDirectory)

	//Config Cors
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowMethods:     enum.CorsAllowMethods,
		AllowHeaders:     enum.CorsAllowHeaders,
		AllowCredentials: true,
		MaxAge:           enum.CorsMaxAge,
	}
	router.Use(cors.New(corsConfig))
	//Set up logger
	zapLog := log.Get(workingDirectory)
	router.Use(logger.RequestLogger(zapLog))

	ctx := context.Background()

	client := db.ConnectDB(ctx, cfg.DatabaseUrl, zapLog)
	//Create Schema
	app := app.Application{
		Config:    cfg,
		EntClient: client,
		Logger:    zapLog,
	}

	if err := client.Schema.Create(ctx); err != nil {
		log.StartUpError(app.Logger, enum.ErrorCreateSchema)
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
		log.StartUpError(app.Logger, enum.ApplicationStartFailed)
	}

}

func LoadRoutes(router *gin.Engine, app *app.Application) {
	routes.ExampleRoutes(router, app)
}
