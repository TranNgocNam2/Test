package certificategrp

import (
	"Backend/business/core/learner/certificate"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	certificate *certificate.Core
}

func New(certificate *certificate.Core) *Handlers {
	return &Handlers{
		certificate: certificate,
	}
}

func (h *Handlers) GetCertificateById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrCertificateIdInvalid)
			return
		}

		certificate, err := h.certificate.GetById(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrCertificateNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, certificate, http.StatusOK, nil)
	}
}

func (h *Handlers) GetCertificates() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageInfo := page.Parse(ctx)

		filter, err := parseFilter(ctx)
		if err != nil {
			filter = certificate.QueryFilter{}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(filterByName, order.ASC)
		}

		certificates, err := h.certificate.Query(ctx, filter, orderBy, pageInfo)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrUserNotFound),
				errors.Is(err, model.ErrSpecNotFound),
				errors.Is(err, model.ErrSubjectNotFound),
				errors.Is(err, model.ErrProgramNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}
		total := h.certificate.Count(ctx, filter)
		result := page.NewPageResponse(certificates, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}

func (h *Handlers) GetSubjectCertificates() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		specId, err := uuid.Parse(ctx.Param("specializationId"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrSpecIDInvalid)
			return
		}

		learnerId := ctx.Query("learnerId")

		certificate, err := h.certificate.GetSubjectCertificates(ctx, specId, &learnerId)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrCertificateNotFound),
				errors.Is(err, model.ErrUserNotFound),
				errors.Is(err, model.ErrSpecNotFound),
				errors.Is(err, model.ErrSubjectNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, certificate, http.StatusOK, nil)
	}
}
