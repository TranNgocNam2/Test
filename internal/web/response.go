package web

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/code"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web"
	"gitlab.com/innovia69420/kit/web/response"
	"net/http"
)

func Respond(ctx *gin.Context, data interface{}, httpStatus int, err error) {
	switch httpStatus {
	case http.StatusOK:
		res := &response.Base{
			ResultCode:    code.Success,
			ResultMessage: message.Success,
			Data:          data,
		}
		ctx.JSON(http.StatusOK, res)
		break

	case http.StatusNotFound:
		web.NotFoundError(ctx, err)
		break

	case http.StatusBadRequest:
		if data != nil {
			web.ValidationError(ctx, err)
			break
		}
		web.BadRequestError(ctx, err)
		break

	case http.StatusInternalServerError:
		web.SystemError(ctx, err)
		break
	}
}
