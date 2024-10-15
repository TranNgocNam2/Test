package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/web"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
)

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
		var request request.NewSubject
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		if err := validateNewSubjectRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		id, err := h.subject.Create(ctx, request)
		if err != nil {
			switch {
			case errors.Is(err, subject.ErrInvalidSkillId),
				errors.Is(err, subject.ErrSkillNotFound),
				errors.Is(err, subject.ErrCodeAlreadyExist):

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

func (h *Handlers) UpdateSubject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		var request request.UpdateSubject
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateSubjectRequest(request); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		status, err := h.subject.GetStatus(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, http.StatusNotFound, err)
		}

		if status == subject.Draft {
			if err = h.subject.UpdateDraft(ctx, request, id); err != nil {
				switch {
				case
					errors.Is(err, subject.ErrSubjectNotFound),
					errors.Is(err, subject.ErrSkillNotFound):

					web.Respond(ctx, nil, http.StatusNotFound, err)
					return

				case
					errors.Is(err, subject.ErrInvalidMaterials),
					errors.Is(err, subject.ErrInvalidSessions):

					web.Respond(ctx, nil, http.StatusBadRequest, err)
					return

				default:
					web.Respond(ctx, nil, http.StatusInternalServerError, err)
					return
				}
			}

			web.Respond(ctx, nil, http.StatusOK, nil)

		} else if status == subject.Published {
			if err = h.subject.UpdatePublished(ctx, request, id); err != nil {
				switch {
				case
					errors.Is(err, subject.ErrSkillNotFound):
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
}

func (h *Handlers) GetSubjectById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		res, err := h.subject.GetById(ctx, id)
		if errors.Is(err, subject.ErrSubjectNotFound) {
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
