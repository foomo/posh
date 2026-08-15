package agent

// Flag is the command-line flag that forces agent mode on, for harnesses that
// set none of the recognized environment variables.
//
// It is exported so callers scanning argv by hand - see
// readline.ExtractFlags, needed where cobra's DisableFlagParsing stops flags
// being populated - do not restate the literal.
const Flag = "--agent"
