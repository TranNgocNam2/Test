package main

import (
	"Backend/api/handlers/schoolgrp"
	"Backend/api/handlers/specializationgrp"
	"Backend/api/handlers/testgrp"
	"Backend/api/handlers/usergrp"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/config"
	"Backend/internal/http"
	"Backend/internal/middleware"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
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

	dbPool, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		fmt.Println(err)
		log.Fatal(message.FailedConnectDatabase)
		return
	}
	defer dbPool.Close()
	db := sqlx.NewDb(stdlib.OpenDBFromPool(dbPool), "pgx")

	a := app.Application{
		Config:  cfg,
		Logger:  log,
		DB:      db,
		Pool:    dbPool,
		Queries: sqlc.New(dbPool),
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
	testgrp.ExampleRoutes(router)

	router.Use(middleware.CheckApiKeyAndRequestID(app.Config.ApiKey))
	usergrp.UserRoutes(router, app)
	schoolgrp.SchoolRoutes(router, app)
	specializationgrp.SpecializationRoutes(router, app)
}
