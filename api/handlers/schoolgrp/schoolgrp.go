package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"errors"
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
		var newSchool request.NewSchool
		if err := web.Decode(ctx, &newSchool); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateCreateSchoolRequest(newSchool); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err := h.school.Create(ctx, toCoreNewSchool(newSchool))
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updatedSchool request.UpdateSchool
		if err := web.Decode(ctx, &updatedSchool); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateSchoolRequest(updatedSchool); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err := h.school.Update(ctx, toCoreUpdateSchool(updatedSchool))
		if err != nil {
			switch {
			case errors.Is(err, school.ErrInvalidID):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case errors.Is(err, school.ErrSchoolNotFound):
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
			switch {
			case errors.Is(err, school.ErrInvalidID):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case errors.Is(err, school.ErrSchoolNotFound):
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

		schoolRes, err := h.school.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, school.ErrSchoolNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toSchoolResponse(schoolRes), http.StatusOK, nil)
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
		result := page.NewPageResponse(toSchoolsResponse(schools), total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSchoolsByDistrict() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		schools, err := h.school.GetSchoolsByDistrictID(ctx)
		if err != nil {
			switch {
			case errors.Is(err, school.ErrInvalidID):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, toSchoolsResponse(schools), http.StatusOK, nil)
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

		web.Respond(ctx, toDistrictsResponse(districts), http.StatusOK, nil)
	}
}
