package slot

import (
	"Backend/business/core/class"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

func (c *Core) UpdateSlot(ctx *gin.Context, id uuid.UUID, updateSlot UpdateSlot) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrCannotUpdateSlotTime
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	slot, err := qtx.GetSlotById(ctx, id)
	if err != nil {
		return model.ErrSlotNotFound
	}

	teacher, err := qtx.GetTeacherById(ctx, updateSlot.TeacherId)
	if err != nil {
		return model.ErrTeacherNotFound
	}

	if slot.StartTime.Equal(updateSlot.StartTime) && slot.EndTime.Equal(updateSlot.EndTime) &&
		teacher.ID == *slot.TeacherID {
		return nil
	}

	if slot.StartTime.Before(time.Now()) {
		return model.ErrSlotAlreadyStarted
	}

	slotConflict, err := qtx.GetConflictingSlotIndexes(ctx, sqlc.GetConflictingSlotIndexesParams{
		ClassID:      slot.ClassID,
		SlotID:       slot.ID,
		NewStartTime: &updateSlot.StartTime,
		NewEndTime:   &updateSlot.EndTime,
	})
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrCannotUpdateSlotTime
	}
	if len(slotConflict) > 0 {
		return fmt.Errorf(model.ErrSlotTimeConflict, slotConflict)
	}

	isTeacherOverlap, err := qtx.CheckTeacherTimeOverlap(ctx, sqlc.CheckTeacherTimeOverlapParams{
		TeacherID: &teacher.ID,
		SlotID:    slot.ID,
		StartTime: &updateSlot.StartTime,
		EndTime:   &updateSlot.EndTime,
	})
	if err != nil || isTeacherOverlap {
		return fmt.Errorf(model.ErrTeacherIsNotAvailable, updateSlot.StartTime.Local().Format(class.TimeLayout),
			updateSlot.EndTime.Local().Format(class.TimeLayout))
	}

	isSlotOrderConflict, err := qtx.CheckSlotOrder(ctx, sqlc.CheckSlotOrderParams{
		SlotID:    slot.ID,
		StartTime: &updateSlot.StartTime,
		EndTime:   &updateSlot.EndTime,
	})
	if err != nil || isSlotOrderConflict {
		return model.ErrSlotOrderTimeConflict
	}

	learnerEmails, _ := qtx.CheckAllLearnersInClassTime(ctx,
		sqlc.CheckAllLearnersInClassTimeParams{
			ClassID:   slot.ClassID,
			EndTime:   &updateSlot.EndTime,
			StartTime: &updateSlot.StartTime,
			SlotID:    slot.ID,
		})
	if len(learnerEmails) > 0 {
		return fmt.Errorf(model.ErrLearnerTimeOverlap, learnerEmails,
			updateSlot.StartTime.Local().Format(class.TimeLayout),
			updateSlot.EndTime.Local().Format(class.TimeLayout))
	}

	err = qtx.UpdateSlotTime(ctx, sqlc.UpdateSlotTimeParams{
		StartTime: &updateSlot.StartTime,
		EndTime:   &updateSlot.EndTime,
		ID:        slot.ID,
		TeacherID: &teacher.ID,
	})
	if err != nil {
		c.logger.Error(err.Error())
		return model.ErrCannotUpdateSlotTime
	}

	tx.Commit(ctx)
	return nil
}
