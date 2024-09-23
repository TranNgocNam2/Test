package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handlers struct {
	school *school.Core
}

func New(school *school.Core) *Handlers {
	return &Handlers{
		school: school,
	}
}

func (h *Handlers) CreateSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var clientNewSchool ClientNewSchool
		if err := web.Decode(ctx, &clientNewSchool); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		newSchool := toCoreNewSchool(clientNewSchool)
		err := h.school.Create(ctx, newSchool)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) DeleteSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err, status := h.school.Delete(ctx)
		if err != nil {
			web.Respond(ctx, nil, status, err)
			return
		}
		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) GetProvinces() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provinces, err := h.school.GetAllProvinces(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		web.Respond(ctx, toClientProvinces(provinces), http.StatusOK, nil)
	}
}

func (h *Handlers) GetDistrictsByProvince() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		districts, err := h.school.GetDistrictsByProvince(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		web.Respond(ctx, toClientDistricts(districts), http.StatusOK, nil)
	}
}
