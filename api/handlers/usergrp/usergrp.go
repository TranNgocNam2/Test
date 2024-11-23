package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
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
		var newUserRequest payload.NewUser
		if err := web.Decode(ctx, &newUserRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
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
				errors.Is(err, model.ErrEmailAlreadyExists),
				errors.Is(err, model.ErrUserAlreadyExist):

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

func (h *Handlers) GetUserById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		userRes, err := h.user.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, userRes, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("id")

		var updateUserRequest payload.UpdateUser
		if err := web.Decode(ctx, &updateUserRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateUserRequest(updateUserRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateUser, err := toCoreUpdateUser(updateUserRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.user.Update(ctx, userID, updateUser)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrPhoneAlreadyExists):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, model.ErrUserNotFound):

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

func (h *Handlers) GetCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRes, err := h.user.GetCurrent(ctx)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, userRes, http.StatusOK, nil)
	}
}

func (h *Handlers) VerifyUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verificationId, err := uuid.Parse(ctx.Param("verificationId"))
		if err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, model.ErrVerificationIdInvalid)
			return
		}
		var verifyUserRequest payload.VerifyLearner
		if err := web.Decode(ctx, &verifyUserRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateVerifyUserRequest(verifyUserRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		verifyUser, err := toCoreVerifyUser(verifyUserRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.user.Verify(ctx, verificationId, verifyUser)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidVerificationInfo),
				errors.Is(err, model.ErrUserCannotBeVerified),
				errors.Is(err, model.ErrInvalidVerificationInfo):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case errors.Is(err, model.ErrUserNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) HandleUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Param("id")

		status, err := h.user.Handle(ctx, userId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}
		data := map[string]string{
			"status": status,
		}
		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) GetVerificationUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo, err := page.Parse(ctx)
		if err != nil {
			pageInfo = page.Page{
				Number: 1,
				Size:   10,
			}
		}

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = user.QueryFilter{
				FullName:   nil,
				SchoolName: nil,
				Status:     nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}
		verificationUsers, err := h.user.GetVerificationUsers(ctx, filter, orderBy, pageInfo)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}
		total := h.user.CountVerificationUsers(ctx, filter)
		result := page.NewPageResponse(verificationUsers, total, pageInfo.Number, pageInfo.Size)
		web.Respond(ctx, result, http.StatusOK, nil)
	}
}
