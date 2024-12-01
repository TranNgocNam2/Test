package payload

type NewUser struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required"`
	Role     *int   `json:"role" validate:"required,gte=0,lte=3"`
}

type UpdateUser struct {
	FullName string `json:"fullName" validate:"required"`
	Photo    string `json:"photo" validate:"required"`
}

type VerifyLearner struct {
	Status int16   `json:"status" validate:"required,gte=1,lte=2"`
	Note   *string `json:"note"`
}

type NewLearner struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullName" validate:"required"`
	Type     *int   `json:"type" validate:"required,gte=0,lte=1"`
	SchoolId string `json:"schoolId" validate:"required"`
}

type UpdateLearner struct {
	Type     *int   `json:"type" validate:"required,gte=0,lte=1"`
	SchoolId string `json:"schoolId" validate:"required"`
}
