package payload

type ClassAccess struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}
