package subject

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/slice"
	"Backend/internal/web/payload"
	"bytes"

	"encoding/json"
	"fmt"
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

func (c *Core) Create(ctx *gin.Context, subject payload.NewSubject) (string, error) {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return "", err
	}

	if _, err := c.queries.IsSubjectCodeExist(ctx, subject.Code); err == nil {
		return "", model.ErrCodeAlreadyExist
	}

	subjectId := uuid.New()
	subjectArgs := sqlc.InsertSubjectParams{
		ID:              subjectId,
		Name:            subject.Name,
		Code:            subject.Code,
		Description:     &subject.Description,
		ImageLink:       &subject.Image,
		Status:          Draft,
		SessionsPerWeek: int16(subject.SessionsPerWeek),
		TimePerSession:  subject.TimePerSession,
		CreatedBy:       staffId,
		CreatedAt:       time.Now().UTC(),
		LearnerType:     subject.LearnerType,
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
		return "", model.ErrInvalidSkillId
	}

	dbSkills, err := qtx.GetSkillsByIds(ctx, skills)
	if err != nil || len(dbSkills) == 0 {
		return "", model.ErrSkillNotFound
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

func (c *Core) UpdateDraft(ctx *gin.Context, s payload.UpdateSubject, id uuid.UUID) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	_, err = c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return model.ErrSubjectNotFound
	}
	totalSessions := len(s.Sessions)
	if *s.Status == Published {
		if len(s.Sessions) == 0 {
			return model.ErrInvalidSessions
		}

		if _, err := c.queries.IsSubjectCodePublished(ctx,
			sqlc.IsSubjectCodePublishedParams{
				Code: s.Code,
				ID:   id,
			}); err == nil {
			return model.ErrCodeAlreadyExist
		}

		if len(s.Transcripts) == 0 {
			return model.ErrInvalidTranscript
		}
		var sum float32
		sum = 0
		for _, transcript := range s.Transcripts {
			sum += transcript.Percentage
		}

		if sum != 100 {
			return model.ErrInvalidTranscriptWeight
		}

		for _, session := range s.Sessions {
			if len(session.Materials) == 0 {
				return model.ErrInvalidMaterials
			}
		}
	}

	skills, err := slice.GetUUIDs(s.Skills)
	if err != nil {
		return model.ErrInvalidSkillId
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	now := time.Now().UTC()

	subParams := sqlc.UpdateSubjectParams{
		Name:            s.Name,
		Code:            s.Code,
		TimePerSession:  s.TimePerSession,
		SessionsPerWeek: int16(s.SessionsPerWeek),
		TotalSessions:   int16(totalSessions),
		MinPassGrade:    &s.MinPassGrade,
		MinAttendance:   &s.MinAttendance,
		Description:     &s.Description,
		Status:          int16(*s.Status),
		ImageLink:       &s.Image,
		UpdatedBy:       &staffId,
		UpdatedAt:       &now,
		LearnerType:     s.LearnerType,
		ID:              id,
	}

	if err := qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if dbSkills, err := qtx.GetSkillsByIds(ctx, skills); err != nil || len(dbSkills) == 0 {
		return model.ErrSkillNotFound
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
				return model.ErrInvalidMaterialType
			}

			materialId, err := uuid.Parse(material.ID)
			if err != nil {
				return fmt.Errorf("Material với id: %s, không đúng định dạng", material.ID)
			}

			data, err := json.Marshal(material.Data)
			if err != nil {
				return model.ErrDataConversion
			}

			param := sqlc.InsertMaterialParams{
				ID:        materialId,
				SessionID: sessionId,
				Type:      material.Type,
				Index:     int32(material.Index),
				IsShared:  material.IsShared,
				Name:      &material.Name,
				Data:      json.RawMessage(data),
			}

			materialParams = append(materialParams, param)
		}

		if _, err := qtx.InsertMaterial(ctx, materialParams); err != nil {
			return err
		}
	}

	if err := qtx.DeleteSubjectTranscripts(ctx, id); err != nil {
		return err
	}

	var transcriptParams []sqlc.InsertTranscriptsParams
	for _, transcript := range s.Transcripts {
		transcriptId, err := uuid.Parse(transcript.Id)
		if err != nil {
			return fmt.Errorf("Transcript với id: %s, không đúng định dạng", transcriptId)
		}

		param := sqlc.InsertTranscriptsParams{
			ID:        transcriptId,
			SubjectID: id,
			Name:      transcript.Name,
			Index:     int32(transcript.Index),
			MinGrade:  float64(transcript.MinGrade),
			Weight:    float64(transcript.Percentage),
		}

		transcriptParams = append(transcriptParams, param)
	}

	if _, err := qtx.InsertTranscripts(ctx, transcriptParams); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) UpdatePublished(ctx *gin.Context, s payload.UpdateSubject, id uuid.UUID) error {
	staffId, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	if _, err := c.queries.IsSubjectCodePublished(ctx,
		sqlc.IsSubjectCodePublishedParams{
			Code: s.Code,
			ID:   id,
		}); err == nil {
		return model.ErrCodeAlreadyExist
	}

	skills, err := slice.GetUUIDs(s.Skills)
	if err != nil {
		return model.ErrInvalidSkillId
	}

	totalSessions := len(s.Sessions)
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)
	now := time.Now().UTC()
	subParams := sqlc.UpdateSubjectParams{
		Name:            s.Name,
		Code:            s.Code,
		TimePerSession:  s.TimePerSession,
		SessionsPerWeek: int16(s.SessionsPerWeek),
		TotalSessions:   int16(totalSessions),
		MinPassGrade:    &s.MinPassGrade,
		MinAttendance:   &s.MinAttendance,
		Description:     &s.Description,
		Status:          int16(*s.Status),
		ImageLink:       &s.Image,
		UpdatedBy:       &staffId,
		UpdatedAt:       &now,
		LearnerType:     s.LearnerType,
		ID:              id,
	}

	if err = qtx.UpdateSubject(ctx, subParams); err != nil {
		return err
	}

	if dbSkills, err := qtx.GetSkillsByIds(ctx, skills); err != nil || len(dbSkills) == 0 {
		return model.ErrSkillNotFound
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
		return nil, model.ErrSubjectNotFound
	}

	totalSessions, err := c.queries.CountSessionsBySubjectId(ctx, id)
	if err != nil {
		totalSessions = 0
	}

	result.ID = subject.ID
	result.Name = subject.Name
	result.Code = subject.Code
	result.TimePerSession = subject.TimePerSession
	if subject.MinAttendance != nil {
		result.MinAttendance = *subject.MinAttendance
	}
	if subject.MinPassGrade != nil {
		result.MinPassGrade = *subject.MinPassGrade
	}
	result.Description = *subject.Description
	result.Image = *subject.ImageLink
	result.LearnerType = *subject.LearnerType
	result.Status = int(subject.Status)
	result.TotalSessions = int(totalSessions)
	result.SessionsPerWeek = int(subject.SessionsPerWeek)
	dbSkills, err := c.queries.GetSkillsBySubjectId(ctx, id)
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

	dbSessions, err := c.queries.GetSessionsBySubjectId(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, dbSession := range dbSessions {
		materials, err := c.queries.GetMaterialsBySessionId(ctx, dbSession.ID)
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
				Type:     dbMaterial.Type,
				Index:    int(dbMaterial.Index),
				IsShared: dbMaterial.IsShared,
				Data:     dbMaterial.Data,
			}

			session.Materials = append(session.Materials, material)
		}

		result.Sessions = append(result.Sessions, session)
	}

	transcripts, err := c.queries.GetTranscriptsBySubjectId(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, dbTranscript := range transcripts {
		transcript := Transcript{
			Id:         dbTranscript.ID,
			Name:       dbTranscript.Name,
			Index:      int(dbTranscript.Index),
			Percentage: float32(dbTranscript.Weight),
			MinGrade:   float32(dbTranscript.MinGrade),
		}

		result.Transcripts = append(result.Transcripts, transcript)
	}

	return &result, nil
}

func (c *Core) GetStatus(ctx *gin.Context, id uuid.UUID) (int, error) {
	subject, err := c.queries.GetSubjectById(ctx, id)
	if err != nil {
		return -1, model.ErrSubjectNotFound
	}

	return int(subject.Status), nil
}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []Subject {
	if err := filter.Validate(); err != nil {
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name, code, image_link, time_per_session, sessions_per_week, description, updated_at, learner_type, status
			FROM subjects`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbSubjects []sqlc.Subject
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbSubjects)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbSubjects == nil {
		return nil
	}

	var subjects []Subject
	for _, dbSubject := range dbSubjects {
		subject := toCoreSubject(dbSubject)
		dbSubjectSkills, err := c.queries.GetSkillsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}
		if dbSubjectSkills != nil {
			for _, skill := range dbSubjectSkills {
				subject.Skills = append(subject.Skills, Skill{
					ID:   skill.ID,
					Name: skill.Name,
				})
			}
		}
		totalSession, err := c.queries.CountSessionsBySubjectId(ctx, dbSubject.ID)
		if err != nil {
			c.logger.Error(err.Error())
			return nil
		}

		subject.TotalSessions = int(totalSession)
		subjects = append(subjects, subject)
	}

	return subjects
}

func (c *Core) Count(ctx *gin.Context, filter QueryFilter) int {
	if err := filter.Validate(); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	data := map[string]interface{}{}

	const q = `SELECT
                        count(1)
               FROM
                        subjects`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}

	if err := pgx.NamedQueryStruct(ctx, c.logger, c.db, buf.String(), data, &count); err != nil {
		c.logger.Error(err.Error())
		return 0
	}

	return count.Count
}

func (c *Core) Delete(ctx *gin.Context, id uuid.UUID) error {
	staffID, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}
	if _, err := c.queries.GetSubjectById(ctx, id); err != nil {
		return model.ErrSubjectNotFound
	}

	if err = c.queries.DeleteSubject(ctx, sqlc.DeleteSubjectParams{
		UpdatedBy: &staffID,
		ID:        id,
	}); err != nil {
		return err
	}
	return nil
}
