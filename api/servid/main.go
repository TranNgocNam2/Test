package main

import (
	"Backend/api/servid/handlers/testgrp"
	"Backend/api/servid/routes"
	"Backend/internal/app"
	"Backend/internal/config"
	"Backend/internal/db"
	"Backend/internal/db/sqlc"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/file"
	"gitlab.com/innovia69420/kit/logger"
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
	//router.Use(log.RequestLogger(zapLog))

	ctx := context.Background()

	dbConn := db.ConnectDB(ctx, cfg.DatabaseUrl, zapLog)

	a := app.Application{
		Config:  cfg,
		Logger:  zapLog,
		Db:      dbConn,
		Queries: sqlc.New(dbConn),
	}

	// Load all routes
	LoadRoutes(router, &a)

	serverAddr := fmt.Sprintf("%s:%d", a.Config.Host, a.Config.Port)
	a.Logger.Info("Server is listening on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		logger.StartUpError(a.Logger, message.FailedStartApplication)
	}

}

func LoadRoutes(router *gin.Engine, app *app.Application) {
	routes.ExampleRoutes(router)
	testgrp.AccountRoutes(router, app)
}
