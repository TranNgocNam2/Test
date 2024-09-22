package schooldb

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SchoolDB struct {
	log *zap.Logger
	db  *sqlx.DB
}

//func (s *SchoolDB) QueryByName(schoolName string, orderBy order.By, page page.Page) (schools *page.Response[school.School], err error) {
//	stmt, err := s.db.Prepare(`SELECT id, name, address, (SELECT COUNT(*) FROM schools WHERE name LIKE %?%) as total FROM schools WHERE name LIKE %?% ORDER BY ? ? LIMIT ? OFFSET ?`)
//	if err != nil {
//		return nil, err
//	}
//	s.db.Na
//	row := stmt.QueryRow(schoolName, schoolName, orderBy.Field, orderBy.Direction, page.Size, (page.Number-1)*page.Size)
//
//	err = row.Scan(&schools)
//	return schools, nil
//}
