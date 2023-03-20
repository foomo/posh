package readline

type Mode string

const (
	ModeArgs           Mode = "args"
	ModeFlags          Mode = "flags"
	ModeAdditionalArgs Mode = "additional"
)
