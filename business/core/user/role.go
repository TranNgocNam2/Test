package user

const (
	LEARNER = iota
	MANAGER
	TEACHER
	ADMIN
)

func GetRoleName(role int16) string {
	switch role {
	case LEARNER:
		return "LEARNER"
	case TEACHER:
		return "TEACHER"
	case MANAGER:
		return "MANAGER"
	case ADMIN:
		return "ADMIN"
	default:
		return "LEARNER"
	}
}

func GetRoleID(role string) int16 {
	switch role {
	case "LEARNER":
		return LEARNER
	case "TEACHER":
		return TEACHER
	case "MANAGER":
		return MANAGER
	case "ADMIN":
		return ADMIN
	default:
		return LEARNER
	}
}
