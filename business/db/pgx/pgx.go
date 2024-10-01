package pgx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

var (
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
	ErrDBNotFound        = errors.New("undefined db")
)

func NamedQuerySlice[T any](ctx *gin.Context, log *zap.Logger, db sqlx.ExtContext, query string, data any, dest *[]T) error {
	return namedQuerySlice(ctx, log, db, query, data, dest, false)
}

func NamedQueryStruct(ctx *gin.Context, log *zap.Logger, db sqlx.ExtContext, query string, data any, dest any) error {
	return namedQueryStruct(ctx, log, db, query, data, dest, false)
}

func namedQuerySlice[T any](ctx *gin.Context, logger *zap.Logger, db sqlx.ExtContext, query string, data any, dest *[]T, withIn bool) (err error) {
	_ = queryString(query, data)

	defer func() {
		if err != nil {
			logger.Info(err.Error())
		}
	}()

	var rows *sqlx.Rows

	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			named, args, err := sqlx.Named(query, data)
			if err != nil {
				return nil, err
			}

			query, args, err := sqlx.In(named, args...)
			if err != nil {
				return nil, err
			}

			query = db.Rebind(query)
			return db.QueryxContext(ctx, query, args...)
		}()

	default:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)
	}

	defer rows.Close()

	var slice []T
	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}
	*dest = slice

	return nil
}

func namedQueryStruct(ctx *gin.Context, log *zap.Logger, db sqlx.ExtContext, query string, data any, dest any, withIn bool) (err error) {
	_ = queryString(query, data)

	defer func() {
		if err != nil {
			log.Info("oopsie daisy")
		}
	}()

	var rows *sqlx.Rows

	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			named, args, err := sqlx.Named(query, data)
			if err != nil {
				return nil, err
			}

			query, args, err := sqlx.In(named, args...)
			if err != nil {
				return nil, err
			}

			query = db.Rebind(query)
			return db.QueryxContext(ctx, query, args...)
		}()

	default:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)
	}

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

func queryString(query string, args any) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " ")
}
