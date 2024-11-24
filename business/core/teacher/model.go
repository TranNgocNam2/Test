package teacher

import "Backend/business/db/sqlc"

type Teacher struct {
	ID       string  `json:"id"`
	FullName string  `json:"fullName"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Photo    *string `json:"photo"`
}

func toCoreTeacher(dbTeacher sqlc.User) Teacher {
	return Teacher{
		ID:       dbTeacher.ID,
		FullName: *dbTeacher.FullName,
		Email:    dbTeacher.Email,
		Phone:    dbTeacher.Phone,
		Photo:    dbTeacher.ProfilePhoto,
	}
}
