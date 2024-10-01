package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.com/innovia69420/kit/web/request"
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
		var request request.NewSchool
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err := validateCreateSchoolRequest(request)

		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.school.Create(ctx, request)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request request.UpdateSchool
		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err := validateUpdateSchoolRequest(request)

		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.school.Update(ctx, request)
		if err != nil {
			switch err {
			case school.ErrInvalidID:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case school.ErrSchoolNotFound:
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

func (h *Handlers) DeleteSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := h.school.Delete(ctx)
		if err != nil {
			switch err {
			case school.ErrInvalidID:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case school.ErrSchoolNotFound:
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

func (h *Handlers) GetSchoolByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		result, err := h.school.GetSchoolByID(ctx, id)
		if err != nil {
			switch err {
			case school.ErrSchoolNotFound:
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toSchoolResponse(*result), http.StatusOK, nil)
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

		schools := h.school.GetSchoolsPaginated(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.school.Count(ctx, filter)

		result := page.NewPageResponse(schools, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSchoolsByDistrict() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schools, err := h.school.GetSchoolsByDistrictID(ctx)
		if err != nil {
			switch err {
			case school.ErrInvalidID:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toWebSchools(schools), http.StatusOK, nil)
	}
}

func (h *Handlers) GetProvinces() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provinces, err := h.school.GetAllProvinces(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, toProvinceResponses(provinces), http.StatusOK, nil)
	}
}

func (h *Handlers) GetDistrictsByProvince() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		districts, err := h.school.GetDistrictsByProvinceID(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, toClientDistricts(districts), http.StatusOK, nil)
	}
}
