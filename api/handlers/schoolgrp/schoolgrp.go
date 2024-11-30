package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/common/model"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		var newSchool payload.NewSchool
		if err := web.Decode(ctx, &newSchool); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateCreateSchoolRequest(newSchool); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		id, err := h.school.Create(ctx, toCoreNewSchool(newSchool))
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		data := map[string]uuid.UUID{
			"id": id,
		}

		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateSchool() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSchoolID)
			return
		}

		var updatedSchool payload.UpdateSchool
		if err := web.Decode(ctx, &updatedSchool); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateUpdateSchoolRequest(updatedSchool); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		err = h.school.Update(ctx, id, toCoreUpdateSchool(updatedSchool))
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSchoolNotFound):
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
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSchoolID)
			return
		}

		err = h.school.Delete(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSchoolNotFound):
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

func (h *Handlers) GetSchoolById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSchoolID)
			return
		}

		schoolRes, err := h.school.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSchoolNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, schoolRes, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSchools() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = school.QueryFilter{
				Name: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}

		schools := h.school.Query(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.school.Count(ctx, filter)
		result := page.NewPageResponse(schools, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSchoolsByDistrict() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidDistrictID)
			return
		}

		schools, err := h.school.GetSchoolsByDistrictId(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, schools, http.StatusOK, nil)
	}
}

func (h *Handlers) GetProvinces() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provinces, err := h.school.GetAllProvinces(ctx)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, provinces, http.StatusOK, nil)
	}
}

func (h *Handlers) GetDistrictsByProvince() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidProvinceID)
			return
		}

		districts, err := h.school.GetDistrictsByProvinceId(ctx, id)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, districts, http.StatusOK, nil)
	}
}
