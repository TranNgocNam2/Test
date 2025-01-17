package web

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	InvalidPayload = errors.New("Dữ liệu không hợp lệ!")
)

type validator interface {
	Validate() error
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
// If the provided value is a struct then it is checked for validation tags.
// If the value implements a validate function, it is executed.
func Decode(c *gin.Context, val any) error {
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return InvalidPayload
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return InvalidPayload
		}
	}

	return nil
}
