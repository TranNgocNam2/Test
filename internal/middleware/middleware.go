package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web"
	"strings"
)

func CheckApiKeyAndRequestID(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.Compare(c.GetHeader(header.XApiKey), apiKey) != 0 {
			web.UnauthorizedError(c, message.InvalidApiKey)
			c.Abort()
			return
		}

		if _, err := uuid.Parse(c.GetHeader(header.XRequestId)); err != nil {
			web.BadRequestError(c, message.InvalidRequestId)
			c.Abort()
			return
		}

		c.Next()
	}
}
