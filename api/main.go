package main

import (
	"Backend/api/cmd/servid/routes"
	"Backend/api/internal/platform/config"
	"Backend/api/internal/platform/db"
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

	//Load all app config
	config.App, _ = config.LoadAllAppConfig(workingDirectory)

	//Config Cors
	corsConfig := cors.Config{
		AllowOrigins:     []string{config.App.CorsOrigin},
		AllowMethods:     enum.CorsAllowMethods,
		AllowHeaders:     enum.CorsAllowHeaders,
		AllowCredentials: true,
		MaxAge:           enum.CorsMaxAge,
	}
	router.Use(cors.New(corsConfig))

	//Set up logger
	logger.Log = log.Get(workingDirectory)
	router.Use(logger.RequestLogger())

	//Connect DB and close connection
	ctx := context.Background()
	db.ConnectDB(ctx, config.App.DatabaseUrl)
	defer db.CloseDB()

	// Load all routes
	LoadRoutes(router)

	serverAddr := fmt.Sprintf("%s:%d", config.App.Host, config.App.Port)
	logger.Log.Info("Server is listening on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.StartUpError(logger.Log, enum.ApplicationStartFailed)
	}

}

func LoadRoutes(router *gin.Engine) {
	routes.ExampleRoutes(router)
}
