package tree

import (
	"context"
	"strings"

	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/util"
	"github.com/spf13/pflag"
)

// CommandInfo describes one command for the `posh agent catalog` catalog. It
// mirrors the shape used by other agent-facing CLIs, so a consumer that can
// read one catalog can read this one.
type CommandInfo struct {
	// FullPath is the command as typed, e.g. "cache clear".
	FullPath    string `json:"full_path"`
	Description string `json:"description"`
	// Dynamic marks a node whose name is resolved at runtime (a cluster name,
	// a cache key, ...). Its FullPath carries a "<placeholder>" instead of a
	// literal name - run the command to discover the actual values.
	Dynamic     bool          `json:"dynamic,omitempty"`
	Arguments   []ArgInfo     `json:"arguments,omitempty"`
	Flags       []FlagInfo    `json:"flags,omitempty"`
	Subcommands []CommandInfo `json:"subcommands,omitempty"`
}

// ArgInfo describes a single positional argument of a command, in the order it
// is expected on the command line.
type ArgInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Optional    bool   `json:"optional,omitempty"`
	// Repeat marks an argument that may be given more than once. Only the last
	// argument of a command can repeat.
	Repeat bool `json:"repeat,omitempty"`
}

// FlagInfo describes a single flag of a command.
type FlagInfo struct {
	Name        string `json:"name"`
	Shorthand   string `json:"shorthand,omitempty"`
	Type        string `json:"type"`
	Default     string `json:"default"`
	Description string `json:"description"`
}

// Usage renders the arguments as a usage string, e.g. "[Key] <Value>...",
// where required arguments are bracketed, optional ones angled and a
// repeatable argument is suffixed with "...".
func (c CommandInfo) Usage() string {
	var ret strings.Builder

	for i, arg := range c.Arguments {
		if i > 0 {
			ret.WriteString(" ")
		}

		if arg.Optional {
			ret.WriteString("<" + arg.Name + ">")
		} else {
			ret.WriteString("[" + arg.Name + "]")
		}

		if arg.Repeat {
			ret.WriteString("...")
		}
	}

	return ret.String()
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

// describe recursively describes a node and its children, accumulating the
// command path as it descends.
func (c *Node) describe(ctx context.Context, r *readline.Readline, prefix string) CommandInfo {
	name := c.Name
	if c.Values != nil {
		if name == "" {
			name = "value"
		}

		name = "<" + name + ">"
	}

	fullPath := util.Pick(prefix != "", prefix+" "+name, name)

	ret := CommandInfo{
		FullPath:    fullPath,
		Description: c.Description,
		Dynamic:     c.Values != nil,
		Arguments:   c.Args.describe(),
		Flags:       c.describeFlags(ctx, r),
	}

	for _, value := range c.Nodes {
		ret.Subcommands = append(ret.Subcommands, value.describe(ctx, r, fullPath))
	}

	return ret
}

// describe describes the positional arguments, preserving the order they are
// expected in.
func (a Args) describe() []ArgInfo {
	if len(a) == 0 {
		return nil
	}

	ret := make([]ArgInfo, 0, len(a))
	for _, arg := range a {
		ret = append(ret, ArgInfo{
			Name:        arg.Name,
			Description: arg.Description,
			Optional:    arg.Optional,
			Repeat:      arg.Repeat,
		})
	}

	return ret
}

// describeFlags enumerates a node's flag definitions. Flags are best effort: a
// node whose Flags callback fails is listed without them rather than failing
// the whole catalog.
func (c *Node) describeFlags(ctx context.Context, r *readline.Readline) []FlagInfo {
	if c.Flags == nil {
		return nil
	}

	fs := readline.NewFlagSets()
	if err := c.Flags(ctx, r, fs); err != nil {
		return nil
	}

	var ret []FlagInfo

	// All merges every set into one pflag.FlagSet, whose VisitAll sorts by
	// name - so the order the sets were populated in does not leak out here.
	fs.All().VisitAll(func(f *pflag.Flag) {
		ret = append(ret, FlagInfo{
			Name:        f.Name,
			Shorthand:   f.Shorthand,
			Type:        f.Value.Type(),
			Default:     f.DefValue,
			Description: f.Usage,
		})
	})

	return ret
}
