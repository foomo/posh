package theme

// Role is a semantic slot in a palette rather than a Catppuccin color name.
//
// Call sites ask for meaning ("this is an error") instead of a color, so
// swapping the palette — in particular flipping between a dark and a light
// flavor — needs no call-site edits. It also keeps the palette data unexported:
// naming Mauve in a public API would promise a distinction the 16-color prompt
// tier cannot keep, since Mocha's Mauve and Lavender collapse to the same
// ANSI slot.
type Role int

const (
	RoleText Role = iota
	RoleMuted
	RolePrimary
	RoleSecondary
	RoleSuccess
	RoleWarning
	RoleError
	RoleInfo
	RoleDebug
)

// roleNames is indexed by Role and must stay in sync with the constants above.
var roleNames = [...]string{
	RoleText:      "text",
	RoleMuted:     "muted",
	RolePrimary:   "primary",
	RoleSecondary: "secondary",
	RoleSuccess:   "success",
	RoleWarning:   "warning",
	RoleError:     "error",
	RoleInfo:      "info",
	RoleDebug:     "debug",
}

// Roles returns every role, in declaration order. Used by tests to assert that
// each palette covers the full set.
func Roles() []Role {
	ret := make([]Role, 0, len(roleNames))
	for i := range roleNames {
		ret = append(ret, Role(i))
	}

	return ret
}

func (r Role) String() string {
	if int(r) < 0 || int(r) >= len(roleNames) {
		return "unknown"
	}

	return roleNames[r]
}
