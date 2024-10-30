package classgrp

import (
	"Backend/business/core/class"
	"Backend/internal/validate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ErrStatusInvalid = errors.New("Trạng thái lớp học không hợp lệ!")
)

const (
	filterByName   = "name"
	filterByCode   = "code"
	filterByStatus = "status"
)

func parseFilter(ctx *gin.Context) (class.QueryFilter, error) {

	var filter class.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	if code := ctx.Query(filterByCode); code != "" {
		filter.WithCode(code)
	}

	if status := ctx.Query(filterByStatus); status != "" {
		statusInt, err := strconv.Atoi(status)
		if err != nil {
			return class.QueryFilter{}, validate.NewFieldsError(filterByStatus, ErrStatusInvalid)
		}

		filter.WithStatus(int16(statusInt))
	}

	return filter, nil
}
