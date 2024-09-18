package http

import (
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/web"
	"strings"
)

func CheckApiKey(c *gin.Context, app *app.Application) {
	if strings.Compare(c.GetHeader(header.XApiKey), app.Config.ApiKey) != 0 {
		web.UnauthorizedError(c)
		return
	}
}
