package readline

type Mode string

const (
	ModeArgs             Mode = ""
	ModeFlags            Mode = "flags"
	ModePassThroughArgs  Mode = "passThrough"
	ModePassThroughFlags Mode = "passThroughFlags"
	ModeAdditionalArgs   Mode = "additional"
)
