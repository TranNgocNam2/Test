package usergrp

import (
	"Backend/api/internal/platform/app"
	"Backend/api/internal/platform/db/ent"
	"Backend/api/internal/platform/db/ent/user"
	"Backend/kit/enum"
	"Backend/kit/web"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func CreateUserHandler(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			web.BadRequest(c, "Dữ liệu không hợp lệ !")
			return
		}

		existingUser, err := app.EntClient.User.Query().Where(user.EmailEQ(req.Email)).Only(c.Request.Context())
		if err == nil && existingUser != nil {
			web.BadRequest(c, "Người dùng đã tồn tại !")
			return
		}

		newUser, err := app.EntClient.User.Create().
			SetName(req.Name).
			SetPassword(req.Password).
			SetEmail(req.Email).
			Save(c.Request.Context())
		if err != nil {
			web.SystemError(c, err)
			return
		}
		response := web.BaseResponse{
			ResultCode:    enum.SuccessCode,
			ResultMessage: enum.SuccessMessage,
			Data: map[string]interface{}{
				"id":       newUser.ID,
				"username": newUser.Name,
				"email":    newUser.Email,
			},
		}
		c.JSON(http.StatusOK, response)
	}
}

func GetAllUsersHandler(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageNumber, err := strconv.Atoi(c.Query("page"))
		if err != nil {
			pageNumber = 1
		}

		pageSize, err := strconv.Atoi(c.Query("size"))
		if err != nil {
			pageSize = 5
		}

		users, err := app.EntClient.User.Query().Order(ent.Asc(user.FieldName)).
			Offset((pageNumber - 1) * pageSize).
			Limit(pageSize).
			All(c.Request.Context())

		response := web.BaseResponse{
			ResultCode:    enum.SuccessCode,
			ResultMessage: enum.SuccessMessage,
			Data:          users,
		}
		c.JSON(http.StatusOK, response)
	}
}
