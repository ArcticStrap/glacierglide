package namespace

const (
	Main int = iota
	Discussion
)

func NamespaceFromNumber(ns int) string {
	switch ns {
	case Main:
		return "main"
	case Discussion:
		return "discussion"
	default:
		return "main"
	}
}
