package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		var newSchoolRequest NewSchoolRequest
		if err := web.Decode(ctx, &newSchoolRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err, statusCode := h.school.Create(ctx, toCoreNewSchool(newSchoolRequest))
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, nil, statusCode, nil)
	}
}

func (h *Handlers) UpdateSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updateSchoolRequest UpdateSchoolRequest
		if err := web.Decode(ctx, &updateSchoolRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err, statusCode := h.school.Update(ctx, toCoreUpdateSchool(updateSchoolRequest))
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, nil, statusCode, nil)
	}
}

func (h *Handlers) DeleteSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err, statusCode := h.school.Delete(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, nil, statusCode, nil)
	}
}

func (h *Handlers) GetSchoolByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
		}

		school, err, statusCode := h.school.GetSchoolByID(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toSchoolResponse(school), statusCode, nil)
	}
}

func (h *Handlers) GetSchoolsByDistrict() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schools, err, statusCode := h.school.GetSchoolsByDistrictID(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toWebSchools(schools), statusCode, nil)
	}
}

func (h *Handlers) GetProvinces() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provinces, err, statusCode := h.school.GetAllProvinces(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toProvinceResponses(provinces), statusCode, nil)
	}
}

func (h *Handlers) GetDistrictsByProvince() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		districts, err, statusCode := h.school.GetDistrictsByProvinceID(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toClientDistricts(districts), statusCode, nil)
	}
}
