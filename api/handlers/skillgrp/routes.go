package skillgrp

import (
	"Backend/business/core/skill"
	"Backend/internal/app"
	"github.com/gin-gonic/gin"
)

func SkillRoutes(router *gin.Engine, app *app.Application) {
	skillCore := skill.NewCore(app)
	handlers := New(skillCore)
	skills := router.Group("/skills")
	{
		skills.POST("", handlers.CreateSkill())
		skills.GET("", handlers.GetSkills())
		skills.PUT("/:id", handlers.UpdateSkill())
		skills.DELETE("/:id", handlers.DeleteSkill())
	}
}
