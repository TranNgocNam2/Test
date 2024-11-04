package payload

type NewUser struct {
	ID    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Role  *int   `json:"role" validate:"required,gte=0,lte=3"`
}

type UpdateUser struct {
	FullName  string   `json:"fullName" validate:"required"`
	Role      *int     `json:"role" validate:"required,gte=0,lte=3"`
	Email     string   `json:"email" validate:"required,email"`
	Phone     string   `json:"phone" validate:"required,startswith=0,len=10"`
	Gender    *int     `json:"gender" validate:"required,gte=0,lte=2"`
	Photo     string   `json:"photo" validate:"required"`
	SchoolID  *string  `json:"schoolID"`
	ImageLink []string `json:"image_links"`
}

type VerifyUser struct {
	Status int `json:"status" validate:"required,gte=1,lte=2"`
}
