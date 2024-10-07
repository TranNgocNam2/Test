package user

const (
	MALE = iota
	FEMALE
	OTHER
)

func GetGenderStr(gender int16) string {
	switch gender {
	case MALE:
		return "MALE"
	case FEMALE:
		return "FEMALE"
	case OTHER:
		return "OTHER"
	default:
		return "OTHER"
	}
}

func GetGenderNum(gender string) int16 {
	switch gender {
	case "MALE":
		return MALE
	case "FEMALE":
		return FEMALE
	case "OTHER":
		return OTHER
	default:
		return OTHER
	}
}
