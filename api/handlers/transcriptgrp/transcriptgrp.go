package transcriptgrp

import (
	"Backend/business/core/transcript"
	"Backend/internal/common/model"
	"Backend/internal/validate"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	transcript *transcript.Core
}

func New(transcript *transcript.Core) *Handlers {
	return &Handlers{
		transcript: transcript,
	}
}

func (h *Handlers) UpdateGrade() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		var request []payload.LearnerTranscript

		if err := web.Decode(ctx, &request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateGradeRequest(request); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		classId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.transcript.ChangeScore(ctx, classId, request)
		if err != nil {
			web.Respond(ctx, nil, http.StatusInternalServerError, err)
			return
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func validateGradeRequest(request []payload.LearnerTranscript) error {
	if err := validate.Check(request); err != nil {
		return err
	}
	return nil
}
