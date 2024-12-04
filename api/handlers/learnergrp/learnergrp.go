package learnergrp

import (
	"Backend/business/core/learner"
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
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case errors.Is(err, model.ErrFailedToAddLearnerToClass):
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
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
			case errors.Is(err, model.ErrInvalidAttendanceCode),
				errors.Is(err, model.ErrSlotNotStarted),
				errors.Is(err, model.ErrSlotEnded):
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

func (h *Handlers) GetLearnerClasses() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("classId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		pageInfo := page.Parse(ctx)

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = learner.QueryFilter{
				FullName:   nil,
				SchoolName: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}

		learners, err := h.learner.GetLearnersInClass(ctx, classId, filter, orderBy, pageInfo.Number, pageInfo.Size)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}
		total := h.learner.CountLearnersInClass(ctx, classId, filter)
		result := page.NewPageResponse(learners, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)

	}
}

func (h *Handlers) GetAttendanceRecords() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		slotId, err := uuid.Parse(ctx.Param("slotId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSlotId)
			return
		}

		pageInfo := page.Parse(ctx)

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = learner.QueryFilter{
				FullName:   nil,
				SchoolName: nil,
				Status:     nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}

		learners, err := h.learner.GetLearnersAttendance(ctx, slotId, filter, orderBy, pageInfo.Number, pageInfo.Size)
		if err != nil {
			web.Respond(ctx, nil, http.StatusUnauthorized, err)
			return
		}

		total := h.learner.CountLearnersAttendance(ctx, slotId, filter)
		result := page.NewPageResponse(learners, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)

	}
}

func (h *Handlers) CreateVerificationInformation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req payload.UpdateVerificationInfo
		if err := web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateLearnerRequest(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		verificationInfo, err := toCoreUpdateLearner(req)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		verificationId, err := h.learner.CreateVerificationInfo(ctx, verificationInfo)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSchoolNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrLearnerAlreadyVerified),
				errors.Is(err, model.ErrVerificationPending):
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
		data := map[string]uuid.UUID{
			"id": verificationId,
		}
		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) CancelVerification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verificationId, err := uuid.Parse(ctx.Param("verificationId"))
		if err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, model.ErrVerificationIdInvalid)
			return
		}
		err = h.learner.CancelVerification(ctx, verificationId)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrVerificationNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrLearnerAlreadyVerified):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser),
				errors.Is(err, model.ErrUnauthorizedFeatureAccess):
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

func (h *Handlers) GetVerificationInfo() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verificationInfo, err := h.learner.GetVerificationsInformation(ctx)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}
		web.Respond(ctx, verificationInfo, http.StatusOK, nil)
	}
}

func (h *Handlers) GetAttendanceReports() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("classId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		reports, err := h.learner.GetAttendanceReports(ctx, classId)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrAttendanceReportsNotFound),
				errors.Is(err, model.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return

			}
		}

		web.Respond(ctx, reports, http.StatusOK, nil)

	}
}
