package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/page"
	"Backend/internal/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/innovia69420/kit/web/request"
	"net/http"
)

var (
	ErrProgramIDInvalid = errors.New("ID khoá học không hợp lệ!")
)

type Handlers struct {
	program *program.Core
}

func New(program *program.Core) *Handlers {
	return &Handlers{
		program: program,
	}
}

func (h *Handlers) CreateProgram() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newProgramRequest request.NewProgram
		if err := web.Decode(ctx, &newProgramRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err := validateNewProgramRequest(newProgramRequest); err != nil {
			fmt.Println(err)

			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		newProgram, err := toCoreNewProgram(newProgramRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		id, err := h.program.Create(ctx, newProgram)
		if err != nil {
			switch {
			case
				errors.Is(err, middleware.ErrInvalidUser):
				web.Respond(ctx, nil, http.StatusUnauthorized, err)
				return
			default:
				web.Respond(ctx, nil, http.StatusInternalServerError, err)
				return
			}
		}

		data := map[string]uuid.UUID{
			"id": id,
		}
		web.Respond(ctx, data, http.StatusOK, nil)
	}
}

func (h *Handlers) UpdateProgram() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, ErrProgramIDInvalid)
			return
		}

		var updateProgramRequest UpdateProgram
		if err = web.Decode(ctx, &updateProgramRequest); err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		if err = validateUpdateProgramRequest(updateProgramRequest); err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, err)
			return
		}

		updateProgram, err := toCoreUpdateProgram(updateProgramRequest)
		if err != nil {
			web.Respond(ctx, nil, http.StatusBadRequest, err)
			return
		}

		err = h.program.Update(ctx, id, updateProgram)
		if err != nil {
			switch {
			case
				errors.Is(err, program.ErrCannotUpdateProgram):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, program.ErrProgramNotFound),
				errors.Is(err, program.ErrSubjectNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
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

func (h *Handlers) DeleteProgram() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			web.Respond(ctx, err, http.StatusBadRequest, ErrProgramIDInvalid)
			return
		}

		err = h.program.Delete(ctx, id)
		if err != nil {
			switch {
			case
				errors.Is(err, program.ErrCannotDeleteProgram):

				web.Respond(ctx, nil, http.StatusBadRequest, err)
				return
			case
				errors.Is(err, program.ErrProgramNotFound):

				web.Respond(ctx, nil, http.StatusNotFound, err)
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

func (h *Handlers) GetPrograms() gin.HandlerFunc {
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
			filter = program.QueryFilter{
				Name: nil,
			}
		}

		orderBy, err := parseOrder(ctx)
		if err != nil {
			orderBy = order.NewBy(orderByStartDate, order.DESC)
		}

		programs := h.program.Query(ctx, filter, orderBy, pageInfo.Number, pageInfo.Size)
		total := h.program.Count(ctx, filter)
		result := page.NewPageResponse(toCoreProgramsResponse(programs), total, pageInfo.Number, pageInfo.Size)

		web.Respond(ctx, result, http.StatusOK, nil)
	}
}
