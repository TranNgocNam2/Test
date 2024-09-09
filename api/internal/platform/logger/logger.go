package logger

import (
	"Backend/kit/enum"
	"Backend/kit/log"
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ResponseWriterWithCapture struct {
	io.Writer
	body *strings.Builder
	gin.ResponseWriter
}

func (r *ResponseWriterWithCapture) Write(b []byte) (int, error) {
	if r.body != nil {
		r.body.Write(b)
	}
	return r.ResponseWriter.Write(b)
}

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		logger = log.FromCtx(context.Request.Context())
		requestId := context.Request.Header.Get(enum.RequestIDHeader)
		if context.Request.Method == http.MethodGet {
			logger.Info(enum.LogStartMessage,
				zap.String(enum.LogRequestId, requestId),
				zap.String(enum.LogMethod, context.Request.Method),
				zap.String(enum.LogUrl, context.Request.RequestURI),
			)
		} else {
			requestBody := ""
			if context.Request.Body != nil && context.Request.Header.Get(enum.HttpContentType) == enum.HttpJson {
				bodyBytes, err := io.ReadAll(context.Request.Body)
				if err == nil {
					requestBody = removeExtraSpacing(string(bodyBytes))
					context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				} else {
					logger.Error(enum.ErrorReadRequestBody, zap.Error(err))
				}
			}

			logger.Info(enum.LogStartMessage,
				zap.String(enum.LogRequestId, requestId),
				zap.String(enum.LogMethod, context.Request.Method),
				zap.String(enum.LogUrl, context.Request.RequestURI),
				zap.String(enum.LogRequestBody, requestBody),
			)
		}

		start := time.Now()

		responseBodyWriter := &ResponseWriterWithCapture{
			body:           &strings.Builder{},
			ResponseWriter: context.Writer,
		}
		context.Writer = responseBodyWriter

		context.Next()

		latency := time.Since(start)

		responseBody := ""
		if responseBodyWriter.body != nil {
			responseBody = responseBodyWriter.body.String()
		}

		logger.Info(enum.LogEndMessage,
			zap.String(enum.LogRequestId, requestId),
			zap.Int(enum.LogStatusCode, context.Writer.Status()),
			zap.String(enum.LogResponseBody, responseBody),
			zap.Duration(enum.LogElapsedTime, latency),
		)
	}
}

func removeExtraSpacing(jsonStr string) string {
	// Remove multiple spaces between key-value pairs
	jsonStr = regexp.MustCompile(`\s+`).ReplaceAllString(jsonStr, " ")

	// Remove leading and trailing spaces within the JSON string
	jsonStr = strings.TrimSpace(jsonStr)

	return jsonStr
}
