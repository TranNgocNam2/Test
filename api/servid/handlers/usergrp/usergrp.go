package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	user *user.Core
}

func New(user *user.Core) *Handlers {
	return &Handlers{
		user: user,
	}
}

func (h *Handlers) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newUserRequest NewUserRequest
		if err := web.Decode(ctx, &newUserRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		newUser, err := toCoreNewUser(newUserRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err, statusCode := h.user.CreateUser(ctx, newUser)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, nil, statusCode, nil)
	}
}

func (h *Handlers) GetUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err, statusCode := h.user.GetUserByID(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toUserResponse(user), statusCode, nil)
	}
}
