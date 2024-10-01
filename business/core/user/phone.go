package user

import (
	"regexp"
)

func IsValidPhoneNumber(phoneNumber string) bool {
	e164Regex := `(0[35789])([0-9]{8})\b`
	re := regexp.MustCompile(e164Regex)
	return re.Find([]byte(phoneNumber)) != nil
}
