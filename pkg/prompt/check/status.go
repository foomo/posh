package check

type Status string

const (
	StatusSuccess = "✅"
	StatusFailure = "❌"
)

func (s Status) String() string {
	return string(s)
}

func StatusFromBool(v bool) Status {
	if v {
		return StatusSuccess
	} else {
		return StatusFailure
	}
}
