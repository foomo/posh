package readline

type Mode string

const (
	ModeArgs             Mode = ""
	ModeFlags            Mode = "flags"
	ModePassThroughFlags Mode = "passThroughFlags"
	ModeAdditionalArgs   Mode = "additional"
)
