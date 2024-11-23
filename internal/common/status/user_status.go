package status

type User int16

const (
	Valid User = iota
	Invalid
)

func GetUserStatus(status int32) string {
	switch User(status) {
	case Invalid:
		return "Invalid"
	default:
		return "Valid"
	}
}
