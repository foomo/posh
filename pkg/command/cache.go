package command

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
)

type Cache struct {
	l     log.Logger
	tree  *tree.Root
	cache cache.Cache
}

// ------------------------------------------------------------------------------------------------
// ~ Constructor
// ------------------------------------------------------------------------------------------------

func NewCache(l log.Logger, cache cache.Cache) *Cache {
	inst := &Cache{
		l:     l,
		cache: cache,
	}
	inst.tree = &tree.Root{
		Name: "cache",
		Nodes: tree.Nodes{
			{
				Name:        "clear",
				Description: "clear all caches",
				Execute:     inst.clear,
			},
			{
				Name:        "list",
				Description: "list all caches",
				Execute:     inst.list,
			},
		},
	}
	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Cache) Name() string {
	return c.tree.Name
}

func (c *Cache) Description() string {
	return "manage the internal cache"
}

func (c *Cache) Complete(ctx context.Context, r *readline.Readline, d prompt.Document) []prompt.Suggest {
	return c.tree.RunCompletion(ctx, r)
}

func (c *Cache) Execute(ctx context.Context, args *readline.Readline) error {
	return c.tree.RunExecution(ctx, args)
}

func (c *Cache) Help() string {
	return `Manage the internal cache.

Usage:
  cache [command]

Available commands:
  list    List all caches
  clear   Clear all caches
`
}

// ------------------------------------------------------------------------------------------------
// ~ Private methods
// ------------------------------------------------------------------------------------------------

func (c *Cache) clear(ctx context.Context, r *readline.Readline) error {
	c.l.Info("clearing cache")
	c.cache.Clear("")
	return nil
}

func (c *Cache) list(ctx context.Context, r *readline.Readline) error {
	// Create a fork of the default table, fill it with data and print it.
	// Data can also be generated and inserted later.
	list := pterm.LeveledList{}
	for ns, value := range c.cache.List() {
		list = append(list, pterm.LeveledListItem{Level: 0, Text: ns})
		for _, k := range value.Keys() {
			list = append(list, pterm.LeveledListItem{Level: 1, Text: k})
			if c.l.Level() == log.LevelTrace {
				list = append(list, pterm.LeveledListItem{Level: 2, Text: fmt.Sprintf("%v", value.Get(k, nil))})
			}
		}
	}
	// Generate tree from LeveledList.
	root := putils.TreeFromLeveledList(list)

	// Render TreePrinter
	return pterm.DefaultTree.WithRoot(root).Render()
}
