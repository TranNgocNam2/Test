package certificategrp

import (
	"Backend/business/core/learner/certificate"
	"Backend/internal/common/model"
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
