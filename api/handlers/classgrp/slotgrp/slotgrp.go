package slotgrp

import (
	"Backend/business/core/class/slot"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"net/http"
)

type Handlers struct {
	slot *slot.Core
}

func New(slot *slot.Core) *Handlers {
	return &Handlers{
		slot: slot,
	}
}

func (h *Handlers) UpdateSlotTime() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrInvalidSlotId)
			return
		}

		var req payload.UpdateSlot
		if err = web.Decode(ctx, &req); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateSlot(req); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateSlot, err := toCoreUpdateSlot(req)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.slot.UpdateSlot(ctx, id, updateSlot)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrSlotNotFound):
				web.Respond(ctx, nil, http.StatusNotFound, err)
				return
			case errors.Is(err, model.ErrCannotUpdateSlotTime):
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			case errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}
