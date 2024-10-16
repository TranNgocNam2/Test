package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/validate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ErrStatusInvalid = errors.New("Trạng thái môn học không hợp lệ!")
)

const (
	filterByName   = "name"
	filterByCode   = "code"
	filterByStatus = "status"
)

func parseFilter(ctx *gin.Context) (subject.QueryFilter, error) {

	var filter subject.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	if code := ctx.Query(filterByCode); code != "" {
		filter.WithCode(code)
	}

	status, err := strconv.Atoi(ctx.DefaultQuery(filterByStatus, "1"))
	if err != nil {
		return subject.QueryFilter{}, validate.NewFieldsError(filterByStatus, ErrStatusInvalid)
	}
	filter.WithStatus(int16(status))

	return filter, nil
}
