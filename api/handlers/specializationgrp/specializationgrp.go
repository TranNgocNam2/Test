package specializationgrp

import (
	"Backend/business/core/specialization"
	"Backend/internal/middleware"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
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
		var newSpecRequest request.NewSpecialization
		if err := web.Decode(ctx, &newSpecRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		if err := validateNewSpecializationRequest(newSpecRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newSpecialization, err := toCoreNewSpecialization(newSpecRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.specialization.Create(ctx, newSpecialization)
		if err != nil {
			switch {
			case
				errors.Is(err, specialization.ErrSkillNotFound),
				errors.Is(err, specialization.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case
				errors.Is(err, specialization.ErrSpecCodeAlreadyExist):
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

func (h *Handlers) GetSpecializationByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		spec, err := h.specialization.GetByID(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, specialization.ErrSpecNotFound),
				errors.Is(err, specialization.ErrSkillNotFound),
				errors.Is(err, specialization.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toResponseSpecialization(spec), http.StatusOK, nil)
	}
}
