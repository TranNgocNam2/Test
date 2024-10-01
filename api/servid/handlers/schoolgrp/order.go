package schoolgrp

import (
	"Backend/business/core/school"
	"Backend/internal/order"
	"Backend/internal/validate"
	"errors"

	"github.com/gin-gonic/gin"
)

func parseOrder(ctx *gin.Context) (order.By, error) {
	const (
		orderByName = "name"
	)

	var orderByFields = map[string]string{
		orderByName: school.OrderByName,
	}

	orderBy, err := order.Parse(ctx, order.NewBy(orderByName, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
