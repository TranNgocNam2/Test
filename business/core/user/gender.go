package user

const (
	MALE    = 1
	FEMALE  = 2
	UNKNOWN = 3
)

func GetGenderStr(gender int16) string {
	switch gender {
	case MALE:
		return "MALE"
	case FEMALE:
		return "FEMALE"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

func GetGenderNum(gender string) int16 {
	switch gender {
	case "MALE":
		return MALE
	case "FEMALE":
		return FEMALE
	case "UNKNOWN":
		return UNKNOWN
	default:
		return UNKNOWN
	}
}
