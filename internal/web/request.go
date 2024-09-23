package web

import (
	"encoding/json"
	"fmt"
	"github.com/dimfeld/httptreemux/v5"
	"github.com/gin-gonic/gin"
	"net/http"
)

type validator interface {
	Validate() error
}

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
func Decode(c *gin.Context, val any) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return fmt.Errorf("unable to decode payload: %w", err)
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("unable to validate payload: %w", err)
		}
	}

	return nil
}
