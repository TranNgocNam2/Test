package routes

import (
	"Backend/kit/enum"
	"Backend/kit/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ExampleRoutes(router *gin.Engine) {
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
}
