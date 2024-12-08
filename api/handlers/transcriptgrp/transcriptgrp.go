package transcriptgrp

import (
	"Backend/business/core/transcript"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/page"
	"Backend/internal/validate"
	"Backend/internal/web"
	"Backend/internal/web/payload"
	"errors"
	"fmt"
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
			fmt.Println(err)
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
			switch {
			case errors.Is(err, model.LearnerNotInClass):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser):

				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) SubmitGrade() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}

		err = h.transcript.SubmitScore(ctx, classId)
		if err != nil {
			switch {
			case
				errors.Is(err, model.ErrClassNotFound),
				errors.Is(err, model.ErrSubjectNotFound),
				errors.Is(err, model.LearnerNotInClass):
				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, middleware.ErrInvalidUser):

				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		web.Respond(ctx, nil, http.StatusOK, nil)
	}
}

func (h *Handlers) GetLearnerTranscripts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		classId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, model.ErrClassIdInvalid)
			return
		}
		pageInfo := page.Parse(ctx)
		filter, err := parseFilter(ctx)
		if err != nil {
			filter = transcript.QueryFilter{
				TranscriptName: nil,
				LearnerId:      nil,
			}
		}

		result := h.transcript.GetLearnerTranscripts(ctx, filter, classId, pageInfo.Number, pageInfo.Size)
		total := h.transcript.Count(ctx, classId, filter)
		results := page.NewPageResponse(result, total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, results, http.StatusOK, nil)
	}
}

func validateGradeRequest(request []payload.LearnerTranscript) error {
	for _, data := range request {
		if err := validate.Check(data); err != nil {
			return err
		}
	}
	return nil
}
