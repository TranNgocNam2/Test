package transcriptgrp

import (
	"Backend/business/core/transcript"

	"github.com/gin-gonic/gin"
)

const (
	filterByName      = "transcriptName"
	filterByLearnerId = "learnerId"
)

func parseFilter(ctx *gin.Context) (transcript.QueryFilter, error) {
	var filter transcript.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	if name := ctx.Query(filterByLearnerId); name != "" {
		filter.WithLearnerId(name)
	}

	return filter, nil
}
