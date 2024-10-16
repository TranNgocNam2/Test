package subject

const (
	Draft = iota
	Published
	Deleted
)

var ValidMaterialTypes = []string{"video", "text", "image", "file", "h1", "h2", "h3", "code"}

func IsTypeValid(input string) bool {
	for _, validType := range ValidMaterialTypes {
		if input == validType {
			return true
		}
	}
	return false
}
