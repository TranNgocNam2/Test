package main

import (
	"Backend/api/cmd/servid/routes"
	"Backend/api/internal/platform/config"
	"Backend/api/internal/platform/logging"
	"Backend/kit/enum"
	"Backend/kit/file"
	"Backend/kit/log"
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

	appCfg, _ := config.LoadAllAppConfig(workingDirectory)

	corsConfig := cors.Config{
		AllowOrigins:     []string{appCfg.CorsOrigin},
		AllowMethods:     enum.CorsAllowMethods,
		AllowHeaders:     enum.CorsAllowHeaders,
		AllowCredentials: true,
		MaxAge:           enum.CorsMaxAge,
	}
	router.Use(cors.New(corsConfig))

	l := log.Get(workingDirectory)
	router.Use(logging.RequestLogger())

	// Load all routes
	LoadRoutes(router)

	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	l.Info("Server is listening on " + serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.StartUpError(l, enum.ApplicationStartFailed)
	}

}

func LoadRoutes(router *gin.Engine) {
	routes.ExampleRoutes(router)
}
