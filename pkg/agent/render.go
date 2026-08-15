package agent

// Render emits a command's result in whichever form the caller is expecting:
// the JSON payload in agent mode, the human-formatted output otherwise.
//
// It exists so commands do not each carry their own `if agent.IsAgentMode()`.
// Both sides are closures, so only the branch actually taken does its work -
// building a payload or walking a cache is not paid for when the other form
// wins.
//
// Payloads are emitted bare by Encode, without a type/schema_version envelope:
// a caller invoked a specific command and already knows what it asked for.
//
//	return agent.Render(
//		func() any { return CacheGet{Namespace: ns, Key: key, Value: value} },
//		func() error { c.l.Info(fmt.Sprintf("%v", value)); return nil },
//	)
func Render(payload func() any, human func() error) error {
	if IsAgentMode() {
		return Encode(payload())
	}

	return human()
}
