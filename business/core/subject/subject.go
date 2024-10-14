package subject

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/middleware"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
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

var (
	ErrSkillNotFound    = errors.New("Kỹ năng không có trong hệ thống!")
	ErrSubjectNotFound  = errors.New("Môn học không có trong hệ thống!")
	ErrCodeAlreadyExist = errors.New("Mã môn đã tồn tại!")
	ErrSkillRequired    = errors.New("Môn học cần có ít nhất một kĩ năng!")
)

func (c *Core) Create(ctx *gin.Context, subject Subject) (string, error) {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return "", err
	}

	if _, err := c.queries.GetSubjectByCode(ctx, subject.Code); err == nil {
		return "", ErrCodeAlreadyExist
	}

	subjectArgs := sqlc.InsertSubjectParams{
		ID:              uuid.New(),
		Name:            subject.Name,
		Description:     subject.Description,
		ImageLink:       subject.Image,
		Status:          Draft,
		TimePerSession:  int16(subject.TimePerSession),
		SessionsPerWeek: int16(subject.SessionPerWeek),
		CreatedBy:       staffId,
		CreatedAt:       time.Now(),
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	id, err := qtx.InsertSubject(ctx, subjectArgs)
	if err != nil {
		return "", err
	}

	_, err = qtx.GetSkillsByIDs(ctx, subject.Skills)
	if err != nil {
		return "", ErrSkillNotFound
	}

	subjectSkillArgs := sqlc.InsertSubjectSkillsParams{
		SubjectID: id,
		SkillIds:  subject.Skills,
	}

	err = qtx.InsertSubjectSkills(ctx, subjectSkillArgs)
	if err != nil {
		return "", err
	}
	tx.Commit(ctx)

	return id.String(), nil
}

func (c *Core) UpdateDraft(ctx *gin.Context, s SubjectDraft) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	now := time.Now()

	subParams := sqlc.UpdateSubjectParams{
		Name:        s.Name,
		Code:        s.Code,
		Description: s.Description,
		Status:      int16(s.Status),
		ImageLink:   s.Image,
		ID:          s.ID,
		UpdatedBy:   &staffId,
		UpdatedAt:   &now,
	}

	if err := qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if _, err = qtx.GetSkillsByIDs(ctx, s.Skills); err == nil {
		return ErrSkillNotFound
	}

	if err := qtx.DeleteSubjectSkills(ctx, s.ID); err != nil {
		return err
	}

	subjectSkillArgs := sqlc.InsertSubjectSkillsParams{
		SubjectID: s.ID,
		SkillIds:  s.Skills,
	}

	err = qtx.InsertSubjectSkills(ctx, subjectSkillArgs)
	if err != nil {
		return err
	}

	for _, session := range s.Sessions {
		sessionParams := sqlc.UpsertSessionParams{
			ID:        session.ID,
			SubjectID: s.ID,
			Index:     int32(session.Index),
			Name:      session.Name,
		}

		if err := qtx.UpsertSession(ctx, sessionParams); err != nil {
			return err
		}

		if err := qtx.DeleteSessionMaterials(ctx, session.ID); err != nil {
			return err
		}

		var materialParams []sqlc.InsertMaterialParams

		for _, material := range session.Materials {
			param := sqlc.InsertMaterialParams{
				ID:        material.ID,
				SessionID: session.ID,
				Index:     int32(material.Index),
				IsShared:  material.IsShared,
				Name:      &material.Name,
				Data:      material.Data,
			}

			materialParams = append(materialParams, param)
		}

		if _, err := qtx.InsertMaterial(ctx, materialParams); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) UpdatePublished(ctx *gin.Context, s SubjectPulished) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)
	now := time.Now()
	subParams := sqlc.UpdateSubjectParams{
		Name:        s.Name,
		Code:        s.Code,
		Description: s.Description,
		Status:      Published,
		ImageLink:   s.Image,
		ID:          s.ID,
		UpdatedBy:   &staffId,
		UpdatedAt:   &now,
	}

	if err = qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if _, err = qtx.GetSkillsByIDs(ctx, s.Skills); err == nil {
		return ErrSkillNotFound
	}

	if err = qtx.DeleteSubjectSkills(ctx, s.ID); err != nil {
		return err
	}

	subjectSkillArgs := sqlc.InsertSubjectSkillsParams{
		SubjectID: s.ID,
		SkillIds:  s.Skills,
	}

	if err = qtx.InsertSubjectSkills(ctx, subjectSkillArgs); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetStatus(ctx *gin.Context, id uuid.UUID) (int, error) {
	subject, err := c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return -1, ErrSubjectNotFound
	}

	return int(subject.Status), nil
}
