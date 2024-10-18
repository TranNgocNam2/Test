package programgrp

import (
	"Backend/business/core/program"
	"Backend/internal/order"
	"Backend/internal/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	InvalidOrderField = errors.New("Thuộc tính cần sắp xếp không hợp lệ: %s!")
)

const (
	orderByName      = "name"
	orderByStartDate = "startDate"
)

func parseOrder(ctx *gin.Context) (order.By, error) {

	var orderByFields = map[string]string{
		orderByName:      program.OrderByName,
		orderByStartDate: program.OrderByStartDate,
	}

	orderBy, err := order.Parse(ctx, order.NewBy(orderByStartDate, order.DESC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, fmt.Errorf(InvalidOrderField.Error(), orderBy.Field))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
