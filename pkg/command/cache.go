package command

import (
	"context"
	"fmt"
	"sort"

	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/samber/lo"
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

func (c *Cache) Complete(ctx context.Context, r *readline.Readline) []goprompt.Suggest {
	return c.tree.Complete(ctx, r)
}

func (c *Cache) Execute(ctx context.Context, r *readline.Readline) error {
	return c.tree.Execute(ctx, r)
}

func (c *Cache) Help(ctx context.Context, r *readline.Readline) string {
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
	cacheList := c.cache.List()
	cacheListKeys := lo.Keys(cacheList)
	sort.Strings(cacheListKeys)
	for _, ns := range cacheListKeys {
		value := cacheList[ns]
		list = append(list, pterm.LeveledListItem{Level: 0, Text: ns})
		keys := value.Keys()
		sort.Strings(keys)
		for _, k := range keys {
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
