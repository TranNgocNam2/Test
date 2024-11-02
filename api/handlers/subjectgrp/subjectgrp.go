package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"net/http"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var ()

type Handlers struct {
	subject *subject.Core
}

func New(subject *subject.Core) *Handlers {
	return &Handlers{
		subject,
	}
}

func (h *Handlers) CreateSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request payload.NewSubject
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewSubjectRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		id, err := h.subject.Create(ctx, request)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrInvalidSkillId),
				errors.Is(err, model.ErrSkillNotFound),
				errors.Is(err, model.ErrCodeAlreadyExist):

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

		data := map[string]string{
			"id": id,
		}

		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSubjectIDInvalid)
			return
		}

		var request payload.UpdateSubject
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateSubjectRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		status, err := h.subject.GetStatus(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, http.StatusNotFound, err)
			return
		}

		if status == subject.Draft {
			if err = h.subject.UpdateDraft(ctx, request, id); err != nil {
				switch {
				case
					errors.Is(err, model.ErrSubjectNotFound),
					errors.Is(err, model.ErrSkillNotFound):

					web.Respond(ctx, nil, http.StatusNotFound, err)
					return

				case
					errors.Is(err, model.ErrCodeAlreadyExist),
					errors.Is(err, model.ErrInvalidMaterials),
					errors.Is(err, model.ErrInvalidTranscript),
					errors.Is(err, model.ErrInvalidTranscriptWeight),
					errors.Is(err, model.ErrInvalidSessions):

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
			return
		}

		if status == subject.Published {
			if err = h.subject.UpdatePublished(ctx, request, id); err != nil {
				switch {
				case
					errors.Is(err, model.ErrSkillNotFound):
					web.Respond(ctx, nil, http.StatusNotFound, err)
					return

				case
					errors.Is(err, model.ErrCodeAlreadyExist):
					web.Respond(ctx, nil, http.StatusBadRequest, err)
					return

				default:
					web.Respond(ctx, nil, http.StatusInternalServerError, err)
					return
				}
			}

			web.Respond(ctx, nil, http.StatusOK, nil)
			return
		}
		web.Respond(ctx, nil, http.StatusInternalServerError, err)
	}
}

func (h *Handlers) GetSubjectById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSubjectIDInvalid)
			return
		}

		res, err := h.subject.GetById(ctx, id)
		if errors.Is(err, model.ErrSubjectNotFound) {
			web.Respond(ctx, nil, http.StatusNotFound, err)
			return
		}

		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, res, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSubjects() gin.HandlerFunc {
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
			filter = subject.QueryFilter{
				Name:   nil,
				Code:   nil,
				Status: subject.Draft,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByCode, order.ASC)
		}

		subjects := h.subject.Query(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.subject.Count(ctx, filter)
		result := page.NewPageResponse(subjects, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		subjectID, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSubjectIDInvalid)
			return
		}

		err = h.subject.Delete(ctx, subjectID)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSubjectNotFound):

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
