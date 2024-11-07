package learnergrp

import (
	"Backend/business/core/learner"
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
	learner *learner.Core
}

func New(learner *learner.Core) *Handlers {
	return &Handlers{
		learner: learner,
	}
}

func (h *Handlers) AddLearnerToClass() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req payload.ClassAccess
		if err := web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewClassAccessRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newClassAccess := toCoreClassAccess(req)

		err := h.learner.JoinClass(ctx, newClassAccess)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrClassStarted),
				errors.Is(err, model.ErrWrongPassword):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
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

func (h *Handlers) AddLearnerToSpecialization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		specializationId, err := uuid.Parse(ctx.Param("specializationId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSpecIDInvalid)
			return
		}

		err = h.learner.JoinSpecialization(ctx, specializationId)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSpecNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrAlreadyJoinedSpecialization):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
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

func (h *Handlers) SubmitAttendance() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("classId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		var req payload.LearnerAttendance
		if err := web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateLearnerAttendanceRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		attendanceSubmission := toCoreSubmitAttendance(req)

		err = h.learner.SubmitAttendance(ctx, classId, attendanceSubmission)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.LearnerNotInClass),
				errors.Is(err, model.ErrSlotNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrInvalidAttendanceCode):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
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
