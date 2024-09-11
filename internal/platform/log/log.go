package log

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gitlab.com/innovia69420/kit/enum/http/header"
	"gitlab.com/innovia69420/kit/enum/message"
	"gitlab.com/innovia69420/kit/logger"
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

func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(context *gin.Context) {
		log = logger.FromCtx(context.Request.Context())
		requestId := context.Request.Header.Get(header.XRequestId)
		if context.Request.Method == http.MethodGet {
			log.Info(message.LogStartMessage,
				zap.String(message.LogRequestId, requestId),
				zap.String(message.LogMethod, context.Request.Method),
				zap.String(message.LogUrl, context.Request.RequestURI),
			)
		} else {
			requestBody := ""
			if context.Request.Body != nil && context.Request.Header.Get(header.ContentType) == header.Json {
				bodyBytes, err := io.ReadAll(context.Request.Body)
				if err == nil {
					requestBody = removeExtraSpacing(string(bodyBytes))
					context.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				} else {
					log.Error(message.FailedReadRequestBody, zap.Error(err))
				}
			}

			log.Info(message.LogStartMessage,
				zap.String(message.LogRequestId, requestId),
				zap.String(message.LogMethod, context.Request.Method),
				zap.String(message.LogUrl, context.Request.RequestURI),
				zap.String(message.LogRequestBody, requestBody),
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

		log.Info(message.LogEndMessage,
			zap.String(message.LogRequestId, requestId),
			zap.Int(message.LogStatusCode, context.Writer.Status()),
			zap.String(message.LogResponseBody, responseBody),
			zap.Duration(message.LogElapsedTime, latency),
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
