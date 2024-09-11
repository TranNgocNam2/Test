package usergrp

import (
	"Backend/internal/platform/app"
	"Backend/internal/platform/db/ent"
	"Backend/internal/platform/db/ent/user"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/code"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web"
	"gitlab.com/innovia69420/kit/web/response"
	"net/http"
	"strconv"
)

func CreateUserHandler(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			web.BadRequestError(c, err.Error())
			return
		}

		existingUser, err := app.EntClient.User.Query().Where(user.EmailEQ(req.Email)).Only(c.Request.Context())
		if err == nil && existingUser != nil {
			web.BadRequestError(c, "Email already exists!")
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
		res := response.Base{
			ResultCode:    code.Success,
			ResultMessage: message.Success,
			Data: map[string]interface{}{
				"id":       newUser.ID,
				"username": newUser.Name,
				"email":    newUser.Email,
			},
		}
		c.JSON(http.StatusOK, res)
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

		res := response.Base{
			ResultCode:    code.Success,
			ResultMessage: message.Success,
			Data:          users,
		}
		c.JSON(http.StatusOK, res)
	}
}
