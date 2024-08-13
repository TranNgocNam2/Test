package web

import (
	"Backend/kit/enum"
	"Backend/kit/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SystemError(context *gin.Context, error error) {
	var systemError = &ErrorResponse{
		ResultCode:    enum.SystemErrorCode,
		ResultMessage: enum.SystemErrorMessage,
	}
	LogRequestError(context, error.Error())
	context.JSON(http.StatusInternalServerError, systemError)
}

func NotFoundError(context *gin.Context, message string) {
	var notFoundError = &ErrorResponse{
		ResultCode:    enum.NotFoundErrorCode,
		ResultMessage: message,
	}
	context.JSON(http.StatusNotFound, notFoundError)
}

func ClientError(context *gin.Context) {
	var clientError = &ErrorResponse{
		ResultCode:    enum.ClientErrorCode,
		ResultMessage: enum.ClientErrorMessage,
	}
	context.JSON(http.StatusBadRequest, clientError)
}

func ResourceError(context *gin.Context) {
	var resourceError = &ErrorResponse{
		ResultCode:    enum.ResourceErrorCode,
		ResultMessage: enum.ResourceErrorMessage,
	}
	context.JSON(http.StatusBadRequest, resourceError)
}

func UnauthorizedError(context *gin.Context) {
	var unAuthorizedError = &ErrorResponse{
		ResultCode:    enum.UnAuthorizedErrorCode,
		ResultMessage: enum.UnAuthorizedErrorMessage,
	}
	context.JSON(http.StatusForbidden, unAuthorizedError)
}

func LogRequestError(context *gin.Context, errMsg string) {
	l := log.FromCtx(context)
	l.Error(errMsg)
}
