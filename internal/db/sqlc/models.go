// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"database/sql"
)

type Account struct {
	ID       string         `json:"id"`
	Username sql.NullString `json:"username"`
	Password sql.NullString `json:"password"`
	Email    sql.NullString `json:"email"`
	Phone    sql.NullString `json:"phone"`
	Address  sql.NullString `json:"address"`
}
