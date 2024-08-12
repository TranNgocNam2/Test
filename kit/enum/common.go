package enum

import "time"

var (
	CorsAllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	CorsAllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Request-ID",
		"Access-Control-Allow-Origin"}
	CorsMaxAge = 12 * time.Hour
)

const (
	EnvironmentFile = ".env"
)
