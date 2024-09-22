package websvc

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/code"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web/response"
	"net/http"
)

func Respond(ctx *gin.Context, data interface{}) {
	res := &response.Base{
		ResultCode:    code.Success,
		ResultMessage: message.Success,
		Data:          data,
	}
	ctx.JSON(http.StatusOK, res)
}
