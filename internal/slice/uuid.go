package slice

import "github.com/google/uuid"

func GetUUIDs(uuidsString []string) ([]uuid.UUID, error) {
	var uuidsSlice []uuid.UUID
	if uuidsString == nil {
		return nil, nil
	}
	for _, str := range uuidsString {
		uuid, err := uuid.Parse(str)
		if err != nil {
			return nil, err
		}
		uuidsSlice = append(uuidsSlice, uuid)
	}
	return uuidsSlice, nil
}
