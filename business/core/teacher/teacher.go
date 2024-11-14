package teacher

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/code"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
	pool    *pgxpool.Pool
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.DB,
		queries: app.Queries,
		logger:  app.Logger,
		pool:    app.Pool,
	}
}

func (c *Core) GenerateAttendanceCode(ctx *gin.Context, slotId uuid.UUID) error {
	teacherId, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return err
	}

	slot, err := c.queries.GetSlotById(ctx, slotId)
	if err != nil {
		return model.ErrSlotNotFound
	}

	if strings.Compare(*slot.TeacherID, teacherId) != 0 {
		return model.ErrTeacherIsNotInSlot
	}

	if slot.StartTime.UTC().After(time.Now().UTC()) {
		return model.ErrSlotNotStarted
	}

	if slot.EndTime.UTC().Before(time.Now().UTC()) {
		return model.ErrSlotEnded
	}

	attendanceCode := code.GenerateAttendance(6)
	fmt.Println(attendanceCode)

	err = c.queries.UpdateAttendanceCode(ctx, sqlc.UpdateAttendanceCodeParams{
		AttendanceCode: &attendanceCode,
		ID:             slot.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
