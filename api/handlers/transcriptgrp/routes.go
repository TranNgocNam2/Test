package transcriptgrp

import (
	"Backend/business/core/transcript"
	"Backend/internal/app"

	"github.com/gin-gonic/gin"
)

func TranscriptRoutes(router *gin.Engine, app *app.Application) {
	transcriptCore := transcript.NewCore(app)
	handlers := New(transcriptCore)

	classes := router.Group("/classes")
	{
		classes.PUT("/:id/transcripts", handlers.UpdateGrade())
		classes.POST("/:id/transcripts/submit", handlers.SubmitGrade())
		classes.GET("/:id/transcripts", handlers.GetLearnerTranscripts())
	}

}
