package assignmentgrp

import (
	"Backend/business/core/class/assignment"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/validate"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	assignment *assignment.Core
}

func New(assignment *assignment.Core) *Handlers {
	return &Handlers{
		assignment: assignment,
	}
}

func (h *Handlers) CreateAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request payload.Assignment
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewAssignmentRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		classId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		id, err := h.assignment.CreateAssignment(ctx, classId, request)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrInvalidDeadlineTime),
				errors.Is(err, model.ErrTimeFormat),
				errors.Is(err, model.ErrDataConversion):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		data := map[string]string{
			"id": id,
		}
		web.Respond(ctx, data, http.StatusOK, nil)

	}
}

func (h *Handlers) UpdateAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request payload.Assignment
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewAssignmentRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		classId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		asmId, err := uuid.Parse(ctx.Param("assignmentId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.assignment.UpdateAssignment(ctx, classId, asmId, request)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrClassNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrInvalidDeadlineTime),
				errors.Is(err, model.ErrTimeFormat),
				errors.Is(err, model.InvalidClassAssignment),
				errors.Is(err, model.ErrDataConversion):

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

func (h *Handlers) DeleteAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		asmId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.assignment.DeleteAssignment(ctx, asmId)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrAssignmentNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrAssignmentDeletion):

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

func (h *Handlers) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		asmId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		result, err := h.assignment.GetById(ctx, asmId)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrAssignmentNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}
		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GradeAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request payload.AssignmentGrade
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateGradeRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		asmId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		learnerId, err := uuid.Parse(ctx.Param("learnerId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.assignment.GradeAssignment(ctx, learnerId, asmId, request)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrLearnerAssignmentNotFound),
				errors.Is(err, model.ErrAssignmentNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.LearnerNotInClass),
				errors.Is(err, model.ErrGradingNotStartedAssignment),
				errors.Is(err, model.ErrInvalidAssignmentSubmision):

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

func (h *Handlers) SubmitAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request payload.LearnerSubmission
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		asmId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}
		err = h.assignment.SubmitAssignment(ctx, asmId, request)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrLearnerAssignmentNotFound),
				errors.Is(err, model.ErrAssignmentNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.LearnerNotInClass),
				errors.Is(err, model.ErrDataConversion),
				errors.Is(err, model.ErrInvalidAssignmentSubmision):

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

func (h *Handlers) GetLearnerAssignment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		asmId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		result, err := h.assignment.GetLearnerAssignment(ctx, asmId)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			case
				errors.Is(err, model.ErrLearnerAssignmentNotFound),
				errors.Is(err, model.ErrAssignmentNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}
		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetAssignments() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

		classId, err := uuid.Parse(ctx.Param("classId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		orderBy := order.NewBy("deadline", order.ASC)

		assignments := h.assignment.Query(ctx, classId, orderBy, pageInfo.Number, pageInfo.Size)

		total := h.assignment.Count(ctx, classId)
		result := page.NewPageResponse(assignments, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func validateNewAssignmentRequest(request payload.Assignment) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}

func validateGradeRequest(request payload.AssignmentGrade) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
