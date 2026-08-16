package agent

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

// TreeNode is the agent-mode representation of a single tree entry.
type TreeNode struct {
	Text     string     `json:"text"`
	Children []TreeNode `json:"children,omitempty"`
}

// TreeJSON is the agent-mode representation of a tree.
type TreeJSON struct {
	Nodes []TreeNode `json:"nodes"`
}

// Tree renders a leveled list, emitting nested JSON in agent mode and a PTerm
// tree otherwise.
//
// Commands should prefer this over calling pterm.DefaultTree directly: PTerm
// draws box-drawing characters an agent cannot parse, and this keeps both paths
// in one call rather than requiring every command to branch.
//
// The flat, level-based pterm.LeveledList is rebuilt into real nesting:
//
//	{"nodes":[{"text":"kubectl","children":[{"text":"clusters"}]}]}
func Tree(list pterm.LeveledList) error {
	return Render(
		func() any { return TreeNodes(list) },
		func() error {
			return pterm.DefaultTree.WithRoot(putils.TreeFromLeveledList(list)).Render()
		},
	)
}

// TreeNodes converts a leveled list into its nested agent-mode representation.
//
// It is exported for commands that need to embed a tree in a larger payload
// rather than emit one on its own; most callers want Tree instead.
//
// Items whose level skips more than one step below their predecessor are
// attached to the nearest available parent, matching how PTerm renders them
// rather than rejecting input a command already considers valid.
func TreeNodes(list pterm.LeveledList) TreeJSON {
	ret := TreeJSON{Nodes: []TreeNode{}}

	// parents[i] is the node currently open at level i; appending to its
	// Children requires a pointer, since TreeNode values are copied.
	var parents []*TreeNode

	for _, item := range list {
		node := TreeNode{Text: item.Text}

		// Clamped into range: a negative level is treated as a root, and one
		// that skips more than a step falls back to the deepest open parent.
		level := min(max(item.Level, 0), len(parents))

		if level == 0 {
			ret.Nodes = append(ret.Nodes, node)
			parents = []*TreeNode{&ret.Nodes[len(ret.Nodes)-1]}

			continue
		}

		parent := parents[level-1]
		parent.Children = append(parent.Children, node)

		parents = append(parents[:level], &parent.Children[len(parent.Children)-1])
	}

	return ret
}
