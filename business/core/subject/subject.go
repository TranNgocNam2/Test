package subject

import (
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/middleware"
	"Backend/internal/slice"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"gitlab.com/innovia69420/kit/web/request"
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
	ErrSkillNotFound       = errors.New("Kỹ năng không có trong hệ thống!")
	ErrSubjectNotFound     = errors.New("Môn học không có trong hệ thống!")
	ErrCodeAlreadyExist    = errors.New("Mã môn đã tồn tại!")
	ErrSkillRequired       = errors.New("Môn học cần có ít nhất một kĩ năng!")
	ErrInvalidSkillId      = errors.New("Skill id không phải định dạng uuid!")
	ErrInvalidSessions     = errors.New("Số lượng session cho môn học không hợp lệ!")
	ErrInvalidMaterials    = errors.New("Buổi học phải có ít nhất 1 nội dung!")
	ErrInvalidMaterialType = errors.New("Material có type không phù hợp!")
)

func (c *Core) Create(ctx *gin.Context, subject request.NewSubject) (string, error) {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return "", err
	}

	if _, err := c.queries.GetSubjectByCode(ctx, subject.Code); err == nil {
		return "", ErrCodeAlreadyExist
	}

	subjectId := uuid.New()
	subjectArgs := sqlc.InsertSubjectParams{
		ID:              subjectId,
		Name:            subject.Name,
		Code:            subject.Code,
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

	skills, err := slice.GetUUIDs(subject.Skills)
	if err != nil {
		return "", ErrInvalidSkillId
	}

	dbSkills, err := qtx.GetSkillsByIDs(ctx, skills)
	if err != nil || len(dbSkills) == 0 {
		return "", ErrSkillNotFound
	}

	var subSkillsParams []sqlc.InsertSubjectSkillParams
	for _, skillId := range skills {
		param := sqlc.InsertSubjectSkillParams{
			ID:        uuid.New(),
			SubjectID: id,
			SkillID:   skillId,
		}

		subSkillsParams = append(subSkillsParams, param)
	}

	if _, err = qtx.InsertSubjectSkill(ctx, subSkillsParams); err != nil {
		return "", err
	}
	tx.Commit(ctx)

	return id.String(), nil
}

func (c *Core) UpdateDraft(ctx *gin.Context, s request.UpdateSubject, id uuid.UUID) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	subject, err := c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return ErrSubjectNotFound
	}

	totalSessions := len(s.Sessions)
	if *s.Status == Published {
		if totalSessions%int(subject.SessionsPerWeek) != 0 || totalSessions == 0 {
			return ErrInvalidSessions
		}

		for _, session := range s.Sessions {
			if len(session.Materials) == 0 {
				return ErrInvalidMaterials
			}
		}
	}

	skills, err := slice.GetUUIDs(s.Skills)
	if err != nil {
		return ErrInvalidSkillId
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
		Status:      int16(*s.Status),
		ImageLink:   s.Image,
		ID:          id,
		UpdatedBy:   &staffId,
		UpdatedAt:   &now,
	}

	if err := qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if dbSkills, err := qtx.GetSkillsByIDs(ctx, skills); err != nil || len(dbSkills) == 0 {
		return ErrSkillNotFound
	}

	if err := qtx.DeleteSubjectSkills(ctx, id); err != nil {
		return err
	}

	var subSkillsParams []sqlc.InsertSubjectSkillParams
	for _, skillId := range skills {
		param := sqlc.InsertSubjectSkillParams{
			ID:        uuid.New(),
			SubjectID: id,
			SkillID:   skillId,
		}

		subSkillsParams = append(subSkillsParams, param)
	}

	if _, err = qtx.InsertSubjectSkill(ctx, subSkillsParams); err != nil {
		return err
	}

	for _, session := range s.Sessions {
		sessionId, err := uuid.Parse(session.ID)
		if err != nil {
			return fmt.Errorf("Session với id: %s, không đúng định dạng", sessionId)
		}
		sessionParams := sqlc.UpsertSessionParams{
			ID:        sessionId,
			SubjectID: id,
			Index:     int32(session.Index),
			Name:      session.Name,
		}

		if err := qtx.UpsertSession(ctx, sessionParams); err != nil {
			return err
		}

		if err := qtx.DeleteSessionMaterials(ctx, sessionId); err != nil {
			return err
		}

		var materialParams []sqlc.InsertMaterialParams

		for _, material := range session.Materials {
			if !IsTypeValid(material.Type) {
				return ErrInvalidMaterialType
			}

			materialId, err := uuid.Parse(material.ID)
			if err != nil {
				return fmt.Errorf("Material với id: %s, không đúng định dạng", material.ID)
			}

			param := sqlc.InsertMaterialParams{
				ID:        materialId,
				SessionID: sessionId,
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

func (c *Core) UpdatePublished(ctx *gin.Context, s request.UpdateSubject, id uuid.UUID) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	skills, err := slice.GetUUIDs(s.Skills)
	if err != nil {
		return ErrInvalidSkillId
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
		ID:          id,
		UpdatedBy:   &staffId,
		UpdatedAt:   &now,
	}

	if err = qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if dbSkills, err := qtx.GetSkillsByIDs(ctx, skills); err != nil || len(dbSkills) == 0 {
		return ErrSkillNotFound
	}

	if err = qtx.DeleteSubjectSkills(ctx, id); err != nil {
		return err
	}

	var subSkillsParams []sqlc.InsertSubjectSkillParams
	for _, skillId := range skills {
		param := sqlc.InsertSubjectSkillParams{
			ID:        uuid.New(),
			SubjectID: id,
			SkillID:   skillId,
		}

		subSkillsParams = append(subSkillsParams, param)
	}

	if _, err = qtx.InsertSubjectSkill(ctx, subSkillsParams); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetById(ctx *gin.Context, id uuid.UUID) (*SubjectDetail, error) {
	var result SubjectDetail
	subject, err := c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return nil, ErrSubjectNotFound
	}

	totalSessions, err := c.queries.CountSessionsBySubjectID(ctx, id)
	if err != nil {
		totalSessions = 0
	}

	result.ID = subject.ID
	result.Name = subject.Name
	result.Code = subject.Code
	result.Description = subject.Description
	result.Image = subject.ImageLink
	result.Status = int(subject.Status)
	result.TotalSessions = int(totalSessions)

	dbSkills, err := c.queries.GetSkillsBySubjectID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, dbSkill := range dbSkills {
		skill := Skill{
			ID:   dbSkill.ID,
			Name: dbSkill.Name,
		}

		result.Skills = append(result.Skills, skill)
	}

	dbSessions, err := c.queries.GetSessionsBySubjectID(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, dbSession := range dbSessions {
		materials, err := c.queries.GetMaterialsBySessionID(ctx, dbSession.ID)
		if err != nil {
			return nil, err
		}

		session := Session{
			ID:    dbSession.ID,
			Name:  dbSession.Name,
			Index: int(dbSession.Index),
		}

		for _, dbMaterial := range materials {
			material := Material{
				ID:       dbMaterial.ID,
				Name:     *dbMaterial.Name,
				Index:    int(dbMaterial.Index),
				IsShared: dbMaterial.IsShared,
				Data:     dbMaterial.Data,
			}

			session.Materials = append(session.Materials, material)
		}

		result.Sessions = append(result.Sessions, session)
	}
	return &result, nil
}

func (c *Core) GetStatus(ctx *gin.Context, id uuid.UUID) (int, error) {
	subject, err := c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return -1, ErrSubjectNotFound
	}

	return int(subject.Status), nil
}
