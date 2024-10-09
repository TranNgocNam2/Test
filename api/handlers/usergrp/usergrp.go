package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
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
		var newUserRequest request.NewUser
		if err := web.Decode(ctx, &newUserRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		if err := validateNewUserRequest(newUserRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newUser, err := toCoreNewUser(newUserRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.user.Create(ctx, newUser)
		if err != nil {
			switch {
			case
				errors.Is(err, user.ErrEmailAlreadyExists),
				errors.Is(err, user.ErrUserAlreadyExist):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) GetUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		userRes, err := h.user.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, user.ErrUserNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toUserResponse(userRes), http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updateUserRequest request.UpdateUser
		if err := web.Decode(ctx, &updateUserRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		updateUser, err := toCoreUpdateUser(updateUserRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.user.Update(ctx, updateUser)
		if err != nil {
			switch {
			case
				errors.Is(err, user.ErrPhoneAlreadyExists),
				errors.Is(err, user.ErrEmailAlreadyExists):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, user.ErrUserNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}
