package providers

const (
	TypeLocal = "local"
	TypeFtp   = "ftp"
)

var typesList = map[string]struct{}{
	TypeLocal: {},
	TypeFtp:   {},
}

func IsValidType(s string) bool {
	_, ok := typesList[s]
	return ok
}
