package routes

import (
	"gitlab.com/innovia69420/kit/enum/code"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web/response"

	"github.com/gin-gonic/gin"
	"net/http"
)

func ExampleRoutes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		res := response.Base{
			ResultCode:    code.Success,
			ResultMessage: message.Success,
			Data: map[string]string{
				"message": "Hello World",
			},
		}
		c.JSON(http.StatusOK, res)
	})
}
