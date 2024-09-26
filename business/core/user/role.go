package user

const (
	LEARNER = 1
	TEACHER = 2
	MANAGER = 3
)

func GetRoleName(role int16) string {
	switch role {
	case LEARNER:
		return "LEARNER"
	case TEACHER:
		return "TEACHER"
	case MANAGER:
		return "MANAGER"
	default:
		return "LEARNER"
	}
}

func IsValidRole(role int) bool {
	switch role {
	case LEARNER:
		return true
	case TEACHER:
		return true
	case MANAGER:
		return true
	default:
		return false
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
	default:
		return LEARNER
	}
}
