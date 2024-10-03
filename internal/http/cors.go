package http

import (
	"gitlab.com/innovia69420/kit/enum/http/header"
	"net/http"
	"time"
)

var (
	AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	AllowHeaders = []string{header.Origin, header.ContentLength, header.ContentType, header.Authorization,
		header.XCsrfToken, header.XRequestId, header.AccessControlAllowOrigin, header.XApiKey}
	CorsMaxAge = 12 * time.Hour
)
