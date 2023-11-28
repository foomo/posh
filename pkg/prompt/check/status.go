package check

type Status string

const (
	StatusNote    = "⬜️️"
	StatusSuccess = "✅"
	StatusFailure = "❌"
)

func (s Status) String() string {
	return string(s)
}
