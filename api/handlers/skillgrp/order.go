package skillgrp

import (
	"Backend/business/core/specialization"
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
	)

	var orderByFields = map[string]string{
		orderByName: specialization.OrderByName,
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
