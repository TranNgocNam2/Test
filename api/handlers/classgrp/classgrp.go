package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/middleware"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"net/http"
)

var (
	ErrClassIDInvalid = errors.New("ID lớp học không hợp lệ!")
)

type Handlers struct {
	class *class.Core
}

func New(class *class.Core) *Handlers {
	return &Handlers{
		class: class,
	}
}

func (h *Handlers) CreateClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newClassRequest request.NewClass
		if err := web.Decode(ctx, &newClassRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewClassRequest(newClassRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newClass, err := toCoreNewClass(newClassRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		id, err := h.class.Create(ctx, newClass)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, class.ErrProgramOrSubjectNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, class.ErrInvalidClassStartTime),
				errors.Is(err, class.ErrInvalidWeekDay):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		data := map[string]uuid.UUID{
			"id": id,
		}
		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}
