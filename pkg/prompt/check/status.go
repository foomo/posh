package check

type Status string

const (
	StatusSuccess = "✅"
	StatusFailure = "❌"
)

func (s Status) String() string {
	return string(s)
}
