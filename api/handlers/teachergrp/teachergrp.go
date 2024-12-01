package teachergrp

import (
	"Backend/business/core/teacher"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/web"
	"Backend/internal/web/payload"
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
		attendanceCode, err := h.teacher.GenerateAttendanceCode(ctx, slotId)
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

		data := map[string]string{"attendanceCode": attendanceCode}

		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) GetTeachersInClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("classId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}
		teachers, err := h.teacher.GetTeachersInClass(ctx, classId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrTeacherNotFound):
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
		web.Respond(ctx, teachers, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateRecord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slotId, err := uuid.Parse(ctx.Param("slotId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSlotId)
			return
		}

		var request payload.UpdateRecord
		if err = web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateRecordRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.teacher.UpdateRecordLink(ctx, slotId, toCoreUpdateRecord(request))
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSlotNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser),
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
