package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrClassIdInvalid = errors.New("Mã lớp học không hợp lệ!")
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
		var newClassRequest NewClass
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
			case
				errors.Is(err, class.ErrProgramNotFound),
				errors.Is(err, class.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, class.ErrInvalidClassStartTime),
				errors.Is(err, class.ErrInvalidWeekDay),
				errors.Is(err, class.ErrClassCodeAlreadyExist):

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

func (h *Handlers) GetClassesByManager() gin.HandlerFunc {
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
			filter = class.QueryFilter{
				Name:   nil,
				Code:   nil,
				Status: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByCode, order.ASC)
		}

		classes := h.class.QueryByManager(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.class.Count(ctx, filter)
		result := page.NewPageResponse(classes, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateClassTeacher() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, ErrClassIdInvalid)
			return
		}

		var updateClassTeacher UpdateClassTeacher
		if err := web.Decode(ctx, &updateClassTeacher); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateClassTeacherRequest(updateClassTeacher); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.class.UpdateClassTeacher(ctx, id, updateClassTeacher.TeacherIds)
		if err != nil {
			switch {
			case
				errors.Is(err, class.ErrClassNotFound),
				errors.Is(err, class.ErrTeacherNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser):
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

func (h *Handlers) GetClassById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, ErrClassIdInvalid)
			return
		}

		classRes, err := h.class.GetByID(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, class.ErrClassNotFound),
				errors.Is(err, class.ErrSubjectNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, classRes, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, ErrClassIdInvalid)
			return
		}

		err = h.class.Delete(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, class.ErrClassNotFound):
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

func (h *Handlers) UpdateClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, ErrClassIdInvalid)
			return
		}

		var updateClassRequest UpdateClass
		if err = web.Decode(ctx, &updateClassRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateClassRequest(updateClassRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateClass, err := toCoreUpdateClass(updateClassRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.class.Update(ctx, id, updateClass)
		if err != nil {
			switch {
			case errors.Is(err, class.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, class.ErrClassCodeAlreadyExist):
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

func (h *Handlers) UpdateClassSlot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, ErrClassIdInvalid)
			return
		}

		var updateSlot UpdateSlot
		if err = web.Decode(ctx, &updateSlot); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateSlotRequest(updateSlot); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateSlots, err := toCoreUpdateSlot(updateSlot)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.class.UpdateSlot(ctx, id, updateSlots)
		if err != nil {
			switch {
			case
				errors.Is(err, class.ErrClassNotFound),
				errors.Is(err, class.ErrSlotNotFound),
				errors.Is(err, class.ErrTeacherNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, class.ErrInvalidSlotStartTime),
				errors.Is(err, class.ErrInvalidSlotEndTime),
				errors.Is(err, class.ErrTeacherNotAvailable),
				errors.Is(err, class.ErrInvalidSlotTime),
				errors.Is(err, class.ErrInvalidSlotCount),
				errors.Is(err, class.ErrTeacherIsNotInClass):

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

func (h *Handlers) CheckTeacherConflict() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkTeacherTime CheckTeacherTime
		if err := web.Decode(ctx, &checkTeacherTime); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
		}

		if err := validateCheckTeacherTimeRequest(checkTeacherTime); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		teacherTime, err := toCoreCheckTeacherTime(checkTeacherTime)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		status := h.class.IsTeacherAvailable(ctx, teacherTime)

		response := map[string]bool{
			"isAvailable": !status,
		}

		web.Respond(ctx, response, http.StatusOK, nil)
	}
}
