package payload

type NewUser struct {
	ID       string `json:"id" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:" " validate:"required"`
	Role     *int   `json:"role" validate:"required,gte=0,lte=3"`
}

type UpdateUser struct {
	FullName string `json:"fullName" validate:"required"`
	Phone    string `json:"phone" validate:"required,startswith=0,len=10"`
	Photo    string `json:"photo" validate:"required"`
}

type VerifyLearner struct {
	Status int16 `json:"status" validate:"required,gte=1,lte=2"`
}
