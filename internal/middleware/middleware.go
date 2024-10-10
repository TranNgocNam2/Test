package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web"
	"strings"
)

var (
	ErrInvalidUser = errors.New("Người dùng không hợp lệ!")
)

func CheckApiKeyAndRequestID(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		if strings.Compare(c.GetHeader(header.XApiKey), apiKey) != 0 {
			err := errors.New(message.InvalidApiKey)
			web.UnauthorizedError(c, err)
			c.Abort()
			return
		}

		if _, err := uuid.Parse(c.GetHeader(header.XRequestId)); err != nil {
			err = errors.New(message.InvalidRequestId)
			web.BadRequestError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
