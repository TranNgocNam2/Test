package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/slice"
	"Backend/internal/web"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/enum/code"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/web/request"
	"gitlab.com/innovia69420/kit/web/response"
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

		skills, err := slice.GetUUIDs(request.Skills)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
		}

		s := subject.Subject{
			Name:           request.Name,
			Code:           request.Code,
			Description:    request.Description,
			Image:          request.Image,
			TimePerSession: request.TimePerSession,
			SessionPerWeek: request.SessionPerWeek,
			Skills:         skills,
		}

		id, err := h.subject.Create(ctx, s)
		if err != nil {
			switch {
			case
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

		skills, err := slice.GetUUIDs(request.Skills)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
		}

		status, err := h.subject.GetStatus(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, http.StatusNotFound, err)
		}

		if status == subject.Draft {
			var sessions []subject.Session
			if len(request.Sessions) != 0 {
				for _, s := range request.Sessions {
					var materials []subject.Material
					if len(s.Materials) != 0 {
						for _, m := range s.Materials {
							id, err = uuid.Parse(m.ID)
							if err != nil {
								web.Respond(ctx, nil, http.StatusBadRequest, fmt.Errorf("Invalid uuid for material: %s", m.ID))
								return
							}

							material := subject.Material{
								ID:       id,
								Name:     m.Name,
								Index:    m.Index,
								IsShared: m.IsShared,
								Data:     m.Data,
							}

							materials = append(materials, material)
						}
					}

					id, err = uuid.Parse(s.ID)
					if err != nil {
						web.Respond(ctx, nil, http.StatusBadRequest, fmt.Errorf("Invalid uuid for session: %s", s.ID))
						return
					}

					session := subject.Session{
						ID:        id,
						Name:      s.Name,
						Index:     s.Index,
						Materials: materials,
					}
					sessions = append(sessions, session)
				}
			}

			s := subject.SubjectDraft{
				ID:          id,
				Name:        request.Name,
				Code:        request.Code,
				Description: request.Description,
				Image:       request.Image,
				Status:      request.Status,
				Skills:      skills,
				Sessions:    sessions,
			}

			if err = h.subject.UpdateDraft(ctx, s); err != nil {
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

		} else if status == subject.Published {
			s := subject.SubjectPulished{
				ID:          id,
				Name:        request.Name,
				Code:        request.Code,
				Description: request.Description,
				Image:       request.Image,
				Skills:      skills,
			}

			if err := h.subject.UpdatePublished(ctx, s); err != nil {
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
