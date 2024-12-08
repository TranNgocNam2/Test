package transcript

import (
	"Backend/business/core/learner/certificate"
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/config"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"Backend/internal/tmplt"
	"Backend/internal/web/payload"
	"bytes"
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jmoiron/sqlx"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type Core struct {
	db      *sqlx.DB
	queries *sqlc.Queries
	logger  *zap.Logger
	pool    *pgxpool.Pool
	config  *config.Config
}

func NewCore(app *app.Application) *Core {
	return &Core{
		db:      app.DB,
		queries: app.Queries,
		logger:  app.Logger,
		pool:    app.Pool,
		config:  app.Config,
	}
}

func (c *Core) ChangeScore(ctx *gin.Context, classId uuid.UUID, req []payload.LearnerTranscript) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return middleware.ErrInvalidUser
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, transcript := range req {
		classLearner, err := qtx.GetClassLearnerByClassAndLearner(ctx, sqlc.GetClassLearnerByClassAndLearnerParams{
			LearnerID: transcript.LearnerId,
			ClassID:   classId,
		})

		if err != nil {
			c.logger.Error("learner with id: % is not in class")
			return model.LearnerNotInClass
		}

		learnerTranscript, err := qtx.GetLearnerTranscript(ctx, sqlc.GetLearnerTranscriptParams{
			ClassLearnerID: classLearner.ID,
			TranscriptID:   transcript.TranscriptId,
		})

		if err != nil {
			c.logger.Error("learner with id: %s does not have this transcript")
			return err
		}

		err = qtx.UpdateLearnerTranscriptGrade(ctx, sqlc.UpdateLearnerTranscriptGradeParams{
			ID:    learnerTranscript.ID,
			Grade: &transcript.Grade,
		})

		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) SubmitScore(ctx *gin.Context, classId uuid.UUID) error {
	_, err := middleware.AuthorizeTeacher(ctx, c.queries)
	if err != nil {
		return middleware.ErrInvalidUser
	}

	class, err := c.queries.GetClassById(ctx, classId)
	if err != nil {
		return model.ErrClassNotFound
	}

	subject, err := c.queries.GetSubjectById(ctx, class.SubjectID)
	if err != nil {
		return model.ErrSubjectNotFound
	}

	classLearners, err := c.queries.GetLearnersByClassId(ctx, class.ID)
	if err != nil {
		return model.CannotGetAllLearners
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := c.queries.WithTx(tx)

	for _, learner := range classLearners {
		transcripts, err := qtx.GetLearnerTranscriptByClassLearnerId(ctx, learner.ClassLearnerID)
		if err != nil {
			return err
		}
		var totalGrade float64
		totalGrade = 0
		pass := true
		for _, transcript := range transcripts {
			if float64(*transcript.Grade) < transcript.MinGrade {
				// Update transcript status
				if err = qtx.UpdateTranscriptStatus(ctx, sqlc.UpdateTranscriptStatusParams{
					ClassLearnerID: learner.ClassLearnerID,
					TranscriptID:   transcript.TranscriptID,
					Status:         2,
				}); err != nil {
					return err
				}
				pass = false
			} else {
				if err = qtx.UpdateTranscriptStatus(ctx, sqlc.UpdateTranscriptStatusParams{
					ClassLearnerID: learner.ClassLearnerID,
					TranscriptID:   transcript.TranscriptID,
					Status:         1,
				}); err != nil {
					return err
				}
				totalGrade = float64(*transcript.Grade) * transcript.Weight
			}
		}

		attendaces, err := qtx.CountAttendace(ctx, learner.ClassLearnerID)
		if err != nil {
			return err
		}

		slots, err := qtx.CountSlotsByClassId(ctx, classId)
		if err != nil {
			return err
		}

		if !pass || totalGrade < float64(*subject.MinPassGrade) || math.Ceil(float64(attendaces)/float64(slots)*100) < float64(*subject.MinAttendance) {
			if err = qtx.UpdateClassStatus(ctx, sqlc.UpdateClassStatusParams{
				ID:     learner.ClassLearnerID,
				Status: 0,
			}); err != nil {
				return err
			}
		} else {
			if err = qtx.UpdateClassStatus(ctx, sqlc.UpdateClassStatusParams{
				ID:     learner.ClassLearnerID,
				Status: 1,
			}); err != nil {
				return err
			}

			certId := uuid.New()

			if err = qtx.CreateSubjectCertificate(ctx, sqlc.CreateSubjectCertificateParams{
				ID:        certId,
				LearnerID: learner.ID,
				SubjectID: &subject.ID,
				Name:      subject.Name,
				Status:    certificate.Valid,
				CreatedAt: time.Now(),
			}); err != nil {
				return err
			}
			err := sendCertiMail(learner.Email,
				*learner.FullName,
				c.config.SendGridApiKey,
				c.config.MAIL_DOMAIN,
				c.config.MailName, subject.Name, certId.String())
			if err != nil {
				return err
			}
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (c *Core) GetLearnerTranscripts(ctx *gin.Context, filter QueryFilter, classId uuid.UUID, pageNumber int, rowsPerPage int) []LearnerTranscriptQuery {
	data := map[string]interface{}{
		"class_id":      classId,
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT DISTINCT cl.learner_id, u.full_name, u.email, lt.transcript_id, t.name, lt.grade, lt.status, t.index
                FROM learner_transcripts lt
                JOIN transcripts t ON lt.transcript_id = t.id
                JOIN class_learners cl ON cl.id = lt.class_learner_id
                JOIN users u ON u.id = cl.learner_id
                WHERE cl.class_id = :class_id`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)
	buf.WriteString(orderByClause(order.NewBy(OrderByIndex, order.ASC)))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var learnerTranscripts []struct {
		LearnerId      string    `db:"learner_id"`
		Name           string    `db:"full_name"`
		Email          string    `db:"email"`
		TranscriptId   uuid.UUID `db:"transcript_id"`
		TranscriptName string    `db:"name"`
		Grade          float32   `db:"grade"`
		Status         int16     `db:"status"`
		Index          int32     `db:"index"`
	}
	if err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &learnerTranscripts); err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if learnerTranscripts == nil {
		return nil
	}

	var result []LearnerTranscriptQuery
	for _, t := range learnerTranscripts {
		t := LearnerTranscriptQuery{
			LearnerId:      t.LearnerId,
			Name:           t.Name,
			Email:          t.Email,
			TranscriptId:   t.TranscriptId,
			TranscriptName: t.TranscriptName,
			Grade:          float64(t.Grade),
			Status:         int32(t.Status),
			Index:          int(t.Index),
		}
		result = append(result, t)
	}

	return result
}

func (c *Core) Count(ctx *gin.Context, classId uuid.UUID, filter QueryFilter) int {
	data := map[string]interface{}{
		"classId": classId,
	}

	const q = `SELECT COUNT(1) as count FROM learner_transcripts lt
                JOIN class_learners cl ON cl.id = lt.class_learner_id
                WHERE cl.class_id = :classId`

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

func sendCertiMail(email string, name string, apiKey string, domain string, fromName string, subjectName string, certiId string) error {
	link := "frontend-innovia.vercel.app/certificate/" + certiId
	data := map[string]interface{}{
		"Name":        name,
		"SubjectName": subjectName,
		"Link":        link,
	}

	htmlTemp, err := template.New("email").Parse(tmplt.SuccessHTML)
	if err != nil {
		return err
	}

	t := template.Must(htmlTemp, err)
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return err
	}
	from := mail.NewEmail(fromName, domain)
	subject := "Congratulations, Your Certificate is Ready!"
	to := mail.NewEmail(name, email)
	html := buf.String()
	message := mail.NewSingleEmail(from, subject, to, "", html)
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil || response.StatusCode != 202 {
		return fmt.Errorf("Gui mail that bai")
	}

	return nil
}
