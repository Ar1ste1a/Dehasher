package files

type FileType int32

const (
	JSON FileType = iota
	XML
	YAML
	TEXT
)

func GetFileType(filetype string) FileType {
	switch filetype {
	case "json":
		return JSON
	case "xml":
		return XML
	case "yaml":
		return YAML
	case "txt":
		return TEXT
	default:
		return JSON
	}
}

func (ft FileType) String() string {
	switch ft {
	case JSON:
		return "json"
	case XML:
		return "xml"
	case YAML:
		return "yaml"
	case TEXT:
		return "txt"
	default:
		return "json"
	}
}

func (ft FileType) Extension() string {
	return "." + ft.String()
}
