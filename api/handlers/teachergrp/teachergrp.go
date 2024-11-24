package teachergrp

import (
	"Backend/business/core/teacher"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	teacher *teacher.Core
}

func New(teacher *teacher.Core) *Handlers {
	return &Handlers{
		teacher: teacher,
	}
}

func (h *Handlers) GenerateAttendanceCode() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slotId, err := uuid.Parse(ctx.Param("slotId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSlotId)
			return
		}
		err = h.teacher.GenerateAttendanceCode(ctx, slotId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSlotNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrSlotNotStarted),
				errors.Is(err, model.ErrSlotEnded):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser),
				errors.Is(err, model.ErrTeacherIsNotInSlot):
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