package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/business/websvc"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/web"
)

type Handlers struct {
	school *school.Core
}

func New(school *school.Core) *Handlers {
	return &Handlers{
		school: school,
	}
}

func (h *Handlers) GetProvinces() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provinces, err := h.school.GetAllProvinces(ctx)
		if err != nil {
			web.BadRequestError(ctx, err.Error())
			return
		}

		websvc.Respond(ctx, toClientProvinces(provinces))
	}
}

func (h *Handlers) GetDistrictsByProvince() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		districts, err := h.school.GetDistrictsByProvince(ctx)
		if err != nil {
			web.BadRequestError(ctx, err.Error())
			return
		}

		websvc.Respond(ctx, toClientDistricts(districts))
	}
}
