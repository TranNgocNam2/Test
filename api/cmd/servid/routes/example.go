package routes

import (
	"Backend/api/internal/platform/app"
	"Backend/kit/enum"
	"Backend/kit/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExampleRoutes(router *gin.Engine, app *app.Application) {
	router.GET("/hello-world", func(c *gin.Context) {
		response := web.BaseResponse{
			ResultCode:    enum.SuccessCode,
			ResultMessage: enum.SuccessMessage,
			Data: map[string]string{
				"message": "Hello World",
			},
		}
		c.JSON(http.StatusOK, response)
	})

	router.GET("/health-check", func(c *gin.Context) {
		err := app.EntClient.Schema.Create(c.Request.Context())
		if err != nil {
			web.SystemError(c, err)
			return
		}
		response := web.BaseResponse{
			ResultCode:    enum.SuccessCode,
			ResultMessage: enum.SuccessMessage,
			Data: map[string]string{
				"message": "OK",
			},
		}

		c.JSON(http.StatusOK, response)
	})

}
