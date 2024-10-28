package program

import (
	"Backend/business/db/sqlc"
	"github.com/google/uuid"
	"time"
)

type NewProgram struct {
	ID          uuid.UUID
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}

type UpdateProgram struct {
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}

type Program struct {
	ID           uuid.UUID
	Name         string
	StartDate    time.Time
	EndDate      time.Time
	TotalClasses int64
}

func toCoreProgram(dbProgram sqlc.Program) Program {
	return Program{
		ID:        dbProgram.ID,
		Name:      dbProgram.Name,
		StartDate: dbProgram.StartDate,
		EndDate:   dbProgram.EndDate,
	}
}
