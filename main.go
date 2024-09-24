package main

import (
	"Backend/api/servid/handlers/schoolgrp"
	"Backend/api/servid/handlers/testgrp"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/config"
	"Backend/internal/http"
	"Backend/internal/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/file"
	"gitlab.com/innovia69420/kit/logger"
	"io"
	"os"
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
	dbConn, err := sqlx.Connect("pgx", cfg.DatabaseUrl)
	if err != nil {
		fmt.Println(err)
		log.Fatal(message.FailedConnectDatabase)
		return
	}

	defer func(dbConn *sqlx.DB) {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(message.FailedCloseDatabase)
		}
	}(dbConn)

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
		log.Fatal(message.FailedStartApplication)
		return
	}

}

func LoadRoutes(router *gin.Engine, app *app.Application) {
	router.Use(middleware.CheckApiKeyAndRequestID(app.Config.ApiKey))

	testgrp.ExampleRoutes(router)
	schoolgrp.SchoolRoutes(router, app)
}
