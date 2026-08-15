package readline

// ExtractFlags consumes a leading run of the given flags off args, invoking
// each one's handler as it goes, and returns the remaining pass-through args.
//
// It exists for commands that set cobra's DisableFlagParsing, where cobra never
// populates the flags and so cannot deliver their values. Registering a flag on
// the root command still matters there - without it cobra rejects the token
// during command resolution - but registration alone does not read it. This
// scan is what actually does.
//
// Stopping at the first unrecognized token is load-bearing: everything from the
// command name onward, including flags meant for the target command (`-n` for
// kubectl), must reach the plugin untouched. Only a leading run is consumed, so
// a later occurrence passes through.
//
//	args = readline.ExtractFlags(args, map[string]func(){
//		"--agent": func() { agent.SetFlag(true) },
//	})
//
// Flag names are matched verbatim, dashes included, so a caller can accept
// whatever spelling it registered.
func ExtractFlags(args []string, flags map[string]func()) []string {
	for len(args) > 0 {
		apply, ok := flags[args[0]]
		if !ok {
			return args
		}

		apply()

		args = args[1:]
	}

	return args
}
