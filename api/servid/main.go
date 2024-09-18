package main

import (
	"Backend/api/servid/handlers/testgrp"
	"Backend/api/servid/routes"
	"Backend/internal/app"
	"Backend/internal/config"
	"Backend/internal/db"
	"Backend/internal/db/sqlc"
	"Backend/internal/http"
	"fmt"
	"github.com/gin-contrib/cors"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/file"
	"gitlab.com/innovia69420/kit/logger"
)

var WorkingDirectory string

func init() {
	var err error
	WorkingDirectory, err = file.WorkingDirectory()
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
	cfg, _ := config.LoadAllAppConfig(WorkingDirectory)

	//Config Cors
	corsConfig := cors.Config{
		AllowOrigins:     []string{cfg.CorsOrigin},
		AllowMethods:     http.AllowMethods,
		AllowHeaders:     http.AllowHeaders,
		AllowCredentials: true,
		MaxAge:           http.CorsMaxAge,
	}
	router.Use(cors.New(corsConfig))
	//Set up log
	log := logger.Get(WorkingDirectory)
	router.Use(logger.RequestLogger(log))

	//ctx := context.Background()

	dbConn, err := db.ConnectDB(cfg.DatabaseUrl, log)
	if err != nil {
		log.Error(message.FailedConnectDatabase)
		return
	}

	a := app.Application{
		Config:  cfg,
		Logger:  log,
		Db:      dbConn,
		Queries: sqlc.New(dbConn),
	}

	// Load all routes
	LoadRoutes(router, &a)

	serverAddr := fmt.Sprintf("%s:%d", a.Config.Host, a.Config.Port)
	log.Info("Server is listening on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Error(message.FailedStartApplication)
		return
	}

}

func LoadRoutes(router *gin.Engine, app *app.Application) {
	routes.ExampleRoutes(router)
	testgrp.AccountRoutes(router, app)
}
