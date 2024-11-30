package payload

type UpdateRecord struct {
	Link string `json:"link" validate:"required"`
}
