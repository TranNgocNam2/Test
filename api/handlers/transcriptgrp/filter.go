package transcriptgrp

import (
	"Backend/business/core/transcript"

	"github.com/gin-gonic/gin"
)

const (
	filterByName = "transcriptName"
)

func parseFilter(ctx *gin.Context) (transcript.QueryFilter, error) {
	var filter transcript.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	return filter, nil
}
