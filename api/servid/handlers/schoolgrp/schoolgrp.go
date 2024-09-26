package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"net/http"

	"github.com/gin-gonic/gin"
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
		school, err, statusCode := h.school.GetSchoolByID(ctx)
		if err != nil {
			web.Respond(ctx, nil, statusCode, err)
			return
		}

		web.Respond(ctx, toSchoolResponse(school), statusCode, nil)
	}
}

func (h *Handlers) GetSchoolPaginated() gin.HandlerFunc {
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
			filter = school.QueryFilter{
				Name: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy("name", order.ASC)
		}

		schools, err, getSchoolStatusCode := h.school.GetSchoolsPaginated(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		if err != nil {
			web.Respond(ctx, nil, getSchoolStatusCode, err)
			return
		}
		total, err, countStatusCode := h.school.Count(ctx, filter)
		if err != nil {
			web.Respond(ctx, nil, countStatusCode, err)
			return
		}

		result := page.NewPageResponse[school.School](schools, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, 200, nil)
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
