package web

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

var (
	InvalidPayload        = errors.New("Dữ liệu không hợp lệ: %s!")
	UnableToDecodePayload = errors.New("Không thể giải mã dữ liệu: %s!")
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
		start := strings.Index(err.Error(), `"`) + 1
		end := strings.LastIndex(err.Error(), `"`)

		fieldName := err.Error()[start:end]
		return fmt.Errorf(UnableToDecodePayload.Error(), fieldName)
	}

	if v, ok := val.(validator); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf(InvalidPayload.Error(), err)
		}
	}

	return nil
}
