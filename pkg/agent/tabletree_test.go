package agent_test

import (
	"testing"

	"github.com/foomo/posh/pkg/agent"
	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableRows_WithHeader(t *testing.T) {
	got := agent.TableRows(pterm.TableData{
		{"Name", "Status"},
		{"prod", "ok"},
		{"dev", "down"},
	}, true)

	require.Len(t, got.Rows, 2)
	assert.Equal(t, map[string]string{"Name": "prod", "Status": "ok"}, got.Rows[0])
	assert.Equal(t, map[string]string{"Name": "dev", "Status": "down"}, got.Rows[1])
}

func TestTableRows_WithoutHeader(t *testing.T) {
	got := agent.TableRows(pterm.TableData{
		{"prod", "ok"},
	}, false)

	require.Len(t, got.Rows, 1)
	assert.Equal(t, []string{"prod", "ok"}, got.Rows[0])
}

func TestTableRows_RaggedRow(t *testing.T) {
	got := agent.TableRows(pterm.TableData{
		{"Name", "Status"},
		{"prod"},
	}, true)

	require.Len(t, got.Rows, 1)
	assert.Equal(t, map[string]string{"Name": "prod", "Status": ""}, got.Rows[0])
}

func TestTableRows_Empty(t *testing.T) {
	got := agent.TableRows(pterm.TableData{}, true)

	assert.NotNil(t, got.Rows, "must encode as [] rather than null")
	assert.Empty(t, got.Rows)
}

func TestTreeNodes_Nesting(t *testing.T) {
	got := agent.TreeNodes(pterm.LeveledList{
		{Level: 0, Text: "alpha"},
		{Level: 1, Text: "one"},
		{Level: 2, Text: "deep"},
		{Level: 1, Text: "two"},
		{Level: 0, Text: "beta"},
	})

	require.Len(t, got.Nodes, 2)

	assert.Equal(t, "alpha", got.Nodes[0].Text)
	require.Len(t, got.Nodes[0].Children, 2)
	assert.Equal(t, "one", got.Nodes[0].Children[0].Text)

	require.Len(t, got.Nodes[0].Children[0].Children, 1)
	assert.Equal(t, "deep", got.Nodes[0].Children[0].Children[0].Text)

	assert.Equal(t, "two", got.Nodes[0].Children[1].Text)
	assert.Empty(t, got.Nodes[0].Children[1].Children)

	assert.Equal(t, "beta", got.Nodes[1].Text)
	assert.Empty(t, got.Nodes[1].Children)
}

// TestTreeNodes_ManyRoots guards against pointers into a slice that append
// reallocates: with enough roots the backing array grows, and a parent pointer
// captured before the growth would write into the discarded array.
func TestTreeNodes_ManyRoots(t *testing.T) {
	var list pterm.LeveledList
	for i := range 64 {
		list = append(list,
			pterm.LeveledListItem{Level: 0, Text: string(rune('a' + i%26))},
			pterm.LeveledListItem{Level: 1, Text: "child"},
		)
	}

	got := agent.TreeNodes(list)

	require.Len(t, got.Nodes, 64)

	for i, node := range got.Nodes {
		require.Len(t, node.Children, 1, "root %d lost its child", i)
		assert.Equal(t, "child", node.Children[0].Text)
	}
}

// TestTreeNodes_ManyChildren is the same hazard one level down.
func TestTreeNodes_ManyChildren(t *testing.T) {
	list := pterm.LeveledList{{Level: 0, Text: "root"}}
	for range 64 {
		list = append(list,
			pterm.LeveledListItem{Level: 1, Text: "child"},
			pterm.LeveledListItem{Level: 2, Text: "grandchild"},
		)
	}

	got := agent.TreeNodes(list)

	require.Len(t, got.Nodes, 1)
	require.Len(t, got.Nodes[0].Children, 64)

	for i, child := range got.Nodes[0].Children {
		require.Len(t, child.Children, 1, "child %d lost its grandchild", i)
	}
}

func TestTreeNodes_SkippedLevel(t *testing.T) {
	got := agent.TreeNodes(pterm.LeveledList{
		{Level: 0, Text: "root"},
		{Level: 3, Text: "jumped"},
	})

	require.Len(t, got.Nodes, 1)
	require.Len(t, got.Nodes[0].Children, 1)
	assert.Equal(t, "jumped", got.Nodes[0].Children[0].Text)
}

func TestTreeNodes_Empty(t *testing.T) {
	got := agent.TreeNodes(pterm.LeveledList{})

	assert.NotNil(t, got.Nodes, "must encode as [] rather than null")
	assert.Empty(t, got.Nodes)
}
