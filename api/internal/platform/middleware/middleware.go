package middleware

import (
	"Backend/kit/enum"
	"Backend/kit/web"
	"fmt"
	"github.com/gin-gonic/gin"
)

func RecoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header(enum.HttpContentType, enum.HttpJson)
				web.SystemError(c, fmt.Errorf("%v", err))
			}
		}()

		c.Next()
	}
}
