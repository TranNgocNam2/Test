package slotgrp

import (
	"Backend/business/core/class/slot"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SlotRoutes(router *gin.Engine, app *app.Application) {
	slotCore := slot.NewCore(app)
	handlers := New(slotCore)
	slots := router.Group("/slots")
	{
		slots.PUT("/:id", handlers.UpdateSlotTime())
	}
}
