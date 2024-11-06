package skillgrp

import (
	"Backend/business/core/skill"
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
	skill *skill.Core
}

func New(skill *skill.Core) *Handlers {
	return &Handlers{
		skill: skill,
	}
}

func (h *Handlers) CreateSkill() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newSkillRequest payload.NewSkill
		if err := web.Decode(ctx, &newSkillRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewSkillRequest(newSkillRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newSkill := toCoreNewSkill(newSkillRequest)

		id, err := h.skill.Create(ctx, newSkill)
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

		resData := map[string]uuid.UUID{
			"id": id,
		}
		web.Respond(ctx, resData, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSkill() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSkillId)
			return
		}
		var updateSkillRequest payload.UpdateSkill
		if err = web.Decode(ctx, &updateSkillRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateSkillRequest(updateSkillRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateSkill := toCoreUpdateSkill(updateSkillRequest)

		err = h.skill.Update(ctx, id, updateSkill)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSkillNotFound):
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

func (h *Handlers) GetSkills() gin.HandlerFunc {
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
			filter = skill.QueryFilter{
				Name: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}

		skills := h.skill.Query(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.skill.Count(ctx, filter)
		result := page.NewPageResponse(skills, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteSkill() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSkillId)
			return
		}

		err = h.skill.Delete(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSkillNotFound):
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
