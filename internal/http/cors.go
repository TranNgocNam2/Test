package http

import "time"

var (
	AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-CSRF-Token", "X-Request-ID",
		"Access-Control-Allow-Origin", "X-API-Key"}
	CorsMaxAge = 12 * time.Hour
)
