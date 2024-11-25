package usergrp

import (
	"Backend/business/core/user"
	"Backend/internal/order"
	"Backend/internal/validate"
	"fmt"
	"github.com/gin-gonic/gin"
)

var (
	InvalidOrderField = "Thuộc tính cần sắp xếp không hợp lệ: %s!"
)

func parseOrder(ctx *gin.Context) (order.By, error) {
	const (
		orderByName      = "name"
		orderByCreatedAt = "createdAt"
	)

	var orderByFields = map[string]string{
		orderByName:      user.OrderByFullName,
		orderByCreatedAt: user.OrderByCreatedAt,
	}

	orderBy, err := order.Parse(ctx, order.NewBy(orderByName, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, fmt.Errorf(InvalidOrderField, orderBy.Field))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
