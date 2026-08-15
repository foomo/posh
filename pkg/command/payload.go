package command

// Payload types for agent mode. They are exported so a provider embedding these
// commands - or writing its own against the same shapes - can unmarshal the
// output rather than re-deriving the JSON structure from a struct literal.
//
// Results are emitted bare, without a type/schema_version envelope: a caller
// invoked a specific command and already knows what it asked for. Log lines
// written by pkg/log.AgentJSON keep their "type" discriminator, which is what
// lets a consumer tell an interleaved log line from the result value.

// EnvVar is a single environment variable. Value is omitted unless the listing
// was requested with values.
type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

// EnvList is the result of `env list`.
type EnvList struct {
	Values []EnvVar `json:"values"`
}

// EnvGet is the result of `env get`.
type EnvGet struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CacheNamespace is a single cache namespace and the keys it holds. Values are
// deliberately excluded - use `cache get` to read one.
type CacheNamespace struct {
	Namespace string   `json:"namespace"`
	Keys      []string `json:"keys"`
}

// CacheList is the result of `cache list`.
type CacheList struct {
	Namespaces []CacheNamespace `json:"namespaces"`
}

// CacheGet is the result of `cache get`. Value is nil when the key is unset.
type CacheGet struct {
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
}

// Version is the result of `posh version`. Commit and BuildTime are only
// populated at debug level, matching the fields the human output prints there.
type Version struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
}
