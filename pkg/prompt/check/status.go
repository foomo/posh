package check

type Status string

const (
	StatusNote    = "note"
	StatusWarning = "warning"
	StatusSuccess = "success"
	StatusFailure = "failure"
)

func (s Status) String() string {
	return string(s)
}
