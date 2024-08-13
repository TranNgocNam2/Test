package routes

import (
	"Backend/api/internal/platform/db"
	"Backend/api/internal/post"
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

	router.GET("/health-check", func(c *gin.Context) {
		err := db.DB.Ping(c.Request.Context())
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

	router.GET("/post", post.GetAllPost(db.Queries))
}
