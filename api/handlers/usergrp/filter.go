package usergrp

import (
	"Backend/business/core/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

const (
	filterByName       = "name"
	filterBySchool     = "school"
	filterByStatus     = "status"
	filterByRole       = "role"
	filterByIsVerified = "isVerified"
)

var (
	InvalidUserStatus   = errors.New("Trạng thái người dùng không hợp lệ!")
	FiltersNotSupported = "Thuộc tính %v không được hỗ trợ!"
	FilterFieldRequired = "Thuộc tính %v là bắt buộc!"
	InvalidFilterData   = "Dữ liệu trong thuộc tính %v không hợp lệ!"
)

func parseFilter(ctx *gin.Context) (user.QueryFilter, error) {

	var filter user.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithFullName(name)
	}

	if isVerified := ctx.Query(filterByIsVerified); isVerified != "" {
		isVerifiedBool, err := strconv.ParseBool(isVerified)
		if err != nil {
			return filter, fmt.Errorf(InvalidFilterData, filterByIsVerified)
		}
		filter.WithIsVerified(isVerifiedBool)
	}

	if school := ctx.Query(filterBySchool); school != "" {
		filter.WithSchoolName(school)
	}

	if statusStr := ctx.Query(filterByStatus); statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err != nil {
			return filter, fmt.Errorf(InvalidFilterData, filterByStatus)
		}
		if status < 0 || status > 1 {
			return filter, InvalidUserStatus
		}
		filter.WithStatus(int16(status))
	}

	if roleStr := ctx.Query(filterByRole); roleStr != "" {
		roleInt, err := strconv.Atoi(roleStr)
		if err != nil {
			return filter, fmt.Errorf(InvalidFilterData, filterByRole)
		}
		filter.WithRole(roleInt)
	}
	return filter, nil

}
