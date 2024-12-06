package certificategrp

import (
	"Backend/business/core/learner/certificate"
	"Backend/internal/validate"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

var (
	ErrStatusInvalid = errors.New("Trạng thái chứng chỉ không hợp lệ!")
)

const (
	filterByName      = "name"
	filterByLearnerId = "learnerId"
	filterByStatus    = "status"
)

func parseFilter(ctx *gin.Context) (certificate.QueryFilter, error) {

	var filter certificate.QueryFilter

	if name := ctx.Query(filterByName); name != "" {
		filter.WithName(name)
	}

	if learnerId := ctx.Query(filterByLearnerId); learnerId != "" {
		filter.WithLearnerId(learnerId)
	}

	if status := ctx.DefaultQuery(filterByStatus, "0"); status != "" {
		statusInt, err := strconv.Atoi(status)
		if err != nil {
			return certificate.QueryFilter{}, validate.NewFieldsError(filterByStatus, ErrStatusInvalid)
		}

		filter.WithStatus(int16(statusInt))
	}

	return filter, nil
}
