package skill

import (
	"Backend/business/db/pgx"
	"Backend/business/db/sqlc"
	"Backend/internal/app"
	"Backend/internal/common/model"
	"Backend/internal/middleware"
	"Backend/internal/order"
	"bytes"
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

func (c *Core) Create(ctx *gin.Context, newSkill NewSkill) (uuid.UUID, error) {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return uuid.Nil, err
	}

	dbSkill := sqlc.CreateSkillParams{
		ID:   newSkill.ID,
		Name: newSkill.Name,
	}

	if err = c.queries.CreateSkill(ctx, dbSkill); err != nil {
		return uuid.Nil, err
	}

	return dbSkill.ID, nil
}

func (c *Core) Update(ctx *gin.Context, id uuid.UUID, updateSkill UpdateSkill) error {
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	dbSkill, err := c.queries.GetSkillById(ctx, id)
	if err != nil {
		return model.ErrSkillNotFound
	}

	dbUpdateSkill := sqlc.UpdateSkillParams{
		ID:   dbSkill.ID,
		Name: updateSkill.Name,
	}

	if err = c.queries.UpdateSkill(ctx, dbUpdateSkill); err != nil {
		return err
	}

	return nil
}

func (c *Core) Query(ctx *gin.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) []Skill {
	if err := filter.Validate(); err != nil {
		return nil
	}

	data := map[string]interface{}{
		"offset":        (pageNumber - 1) * rowsPerPage,
		"rows_per_page": rowsPerPage,
	}

	const q = `SELECT
						id, name
			FROM skills`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	buf.WriteString(orderByClause(orderBy))
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")
	c.logger.Info(buf.String())

	var dbSkills []sqlc.Skill
	err := pgx.NamedQuerySlice(ctx, c.logger, c.db, buf.String(), data, &dbSkills)
	if err != nil {
		c.logger.Error(err.Error())
		return nil
	}

	if dbSkills == nil {
		return nil
	}

	var skills []Skill
	for _, dbSkill := range dbSkills {
		skills = append(skills, toCoreSkill(dbSkill))
	}

	return skills
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
                        skills`

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
	_, err := middleware.AuthorizeStaff(ctx, c.queries)
	if err != nil {
		return err
	}

	if _, err = c.queries.GetSkillById(ctx, id); err != nil {
		return model.ErrSkillNotFound
	}

	if err = c.queries.DeleteSkill(ctx, id); err != nil {
		return err
	}

	return nil
}
