package subjectgrp

import (
	"Backend/business/core/subject"
	"Backend/internal/order"
	"Backend/internal/validate"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	InvalidOrderField = errors.New("Thuộc tính cần sắp xếp không hợp lệ: %s!")
)

func parseOrder(ctx *gin.Context) (order.By, error) {
	const (
		orderByName = "name"
		orderByCode = "code"
	)

	var orderByFields = map[string]string{
		orderByName: subject.OrderByName,
		orderByCode: subject.OrderByCode,
	}

	orderBy, err := order.Parse(ctx, order.NewBy(orderByName, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, fmt.Errorf(InvalidOrderField.Error(), orderBy.Field))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
