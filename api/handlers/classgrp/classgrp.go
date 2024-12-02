package classgrp

import (
	"Backend/business/core/class"
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
	class *class.Core
}

func New(class *class.Core) *Handlers {
	return &Handlers{
		class: class,
	}
}

func (h *Handlers) CreateClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newClassRequest payload.NewClass
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
				errors.Is(err, model.ErrProgramNotFound),
				errors.Is(err, model.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrInvalidClassStartTime),
				errors.Is(err, model.ErrInvalidWeekDay),
				errors.Is(err, model.ErrClassCodeAlreadyExist),
				errors.Is(err, model.ErrInvalidSessionCount):

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

func (h *Handlers) ImportLearners() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.ImportLearners
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateImportLearnersRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		learners, err := toCoreImportLearners(req)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.class.ImportLearners(ctx, id, learners)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.CannotGetAllLearners):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrCannotImportLearners):
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) GetClassesByManager() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

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

		classes, err := h.class.QueryByManager(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}
		total := h.class.Count(ctx, filter)
		result := page.NewPageResponse(classes, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetClassesByLearner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classes, err := h.class.QueryByLearner(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}
		web.Respond(ctx, classes, http.StatusOK, nil)
	}
}

func (h *Handlers) GetClassesByTeacher() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

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

		classes, err := h.class.QueryByTeacher(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}

		total := h.class.CountByTeacher(ctx, filter)
		result := page.NewPageResponse(classes, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetClassById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		classRes, err := h.class.GetByID(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrSubjectNotFound):
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

		web.Respond(ctx, classRes, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.class.Delete(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound):
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
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var updateClassRequest payload.UpdateClass
		if err = web.Decode(ctx, &updateClassRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateClassRequest(updateClassRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateClass := toCoreUpdateClass(updateClassRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.class.Update(ctx, id, updateClass)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, model.ErrClassCodeAlreadyExist):
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

func (h *Handlers) UpdateMeetingLink() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.UpdateMeetingLink
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateMeetingLinkRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.class.UpdateMeetingLink(ctx, id, toCoreUpdateMeetingLink(req))
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrClassNotCompleted),
				errors.Is(err, model.ErrClassIsEnded):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser),
				errors.Is(err, model.ErrTeacherIsNotInClass):
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

func (h *Handlers) UpdateClassSlot() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.UpdateSlots
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateSlotRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateSlots, err := toCoreUpdateSlots(req)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.class.UpdateSlots(ctx, id, updateSlots, *req.Status)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrSlotNotFound),
				errors.Is(err, model.ErrTeacherNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, model.ErrInvalidSlotStartTime),
				errors.Is(err, model.ErrInvalidSlotEndTime),
				errors.Is(err, model.ErrTeacherNotAvailable),
				errors.Is(err, model.ErrInvalidSlotTime),
				errors.Is(err, model.ErrInvalidSlotCount):

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

func (h *Handlers) CheckTeacherAvailable() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var checkTeacherTime payload.CheckTeacherTime
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

		isAvailable, err := h.class.IsTeacherAvailable(ctx, teacherTime)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrTeacherNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		res := map[string]bool{
			"isAvailable": isAvailable,
		}

		web.Respond(ctx, res, http.StatusOK, nil)
	}
}

func (h *Handlers) AddLearner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.AddLearner
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateAddLearnerRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.class.AddLearner(ctx, id, toCoreAddLearner(req))
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrLearnerNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrFailedToAddLearnerToClass):
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) RemoveLearner() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.RemoveLearner
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateRemoveLearnerRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.class.RemoveLearner(ctx, id, toCoreRemoveLearner(req))
		if err != nil {
			switch {
			case errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrLearnerNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrLearnerNotInClass),
				errors.Is(err, model.ErrClassStarted):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
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
