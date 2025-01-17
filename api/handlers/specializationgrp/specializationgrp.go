package specializationgrp

import (
	"Backend/business/core/specialization"
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
	specialization *specialization.Core
}

func New(specialization *specialization.Core) *Handlers {
	return &Handlers{
		specialization: specialization,
	}
}

func (h *Handlers) CreateSpecialization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newSpecRequest payload.NewSpecialization
		if err := web.Decode(ctx, &newSpecRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewSpecializationRequest(newSpecRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newSpec := toCoreNewSpecialization(newSpecRequest)

		id, err := h.specialization.Create(ctx, newSpec)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSkillNotFound),
				errors.Is(err, model.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrSpecCodeAlreadyExist):
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

		resData := map[string]uuid.UUID{
			"id": id,
		}
		web.Respond(ctx, resData, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSpecialization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		specID, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSpecIDInvalid)
			return
		}

		var updateSpecRequest payload.UpdateSpecialization
		if err := web.Decode(ctx, &updateSpecRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateSpecializationRequest(updateSpecRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateSpec, err := toCoreUpdatedSpecialization(updateSpecRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.specialization.Update(ctx, specID, updateSpec)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSubjectNotFound),
				errors.Is(err, model.ErrSpecNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, model.ErrSpecCodeAlreadyExist):
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

func (h *Handlers) DeleteSpecialization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		specID, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSpecIDInvalid)
			return
		}

		err = h.specialization.Delete(ctx, specID)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrSpecNotFound):

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

func (h *Handlers) GetSpecializationById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		specialization, err := h.specialization.GetById(ctx, id)
		if err != nil {
			switch {
			//case
			//	errors.Is(err, middleware.ErrInvalidUser):
			//	web.Respond(ctx, nil, http.StatusUnauthorized, err)
			//	return
			case
				errors.Is(err, model.ErrSpecNotFound),
				errors.Is(err, model.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, specialization, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSpecializations() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = specialization.QueryFilter{
				Name:   nil,
				Code:   nil,
				Status: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(specialization.OrderByName, order.ASC)
		}

		specializations := h.specialization.Query(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.specialization.Count(ctx, filter)
		result := page.NewPageResponse(specializations, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}
