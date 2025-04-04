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
	"github.com/foomo/posh/pkg/util/suggests"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/samber/lo"
)

type Cache struct {
	l     log.Logger
	tree  tree.Root
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
	inst.tree = tree.New(&tree.Node{
		Name:        "cache",
		Description: "Manage the internal cache",
		Nodes: tree.Nodes{
			{
				Name:        "clear",
				Description: "Clear caches",
				Args: tree.Args{
					{
						Name:        "Namespace",
						Description: "Name of namespace to clear.",
						Repeat:      true,
						Optional:    true,
						Suggest: func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
							return suggests.List(lo.Keys(inst.cache.List()))
						},
					},
				},
				Execute: inst.clear,
			},
			{
				Name:        "list",
				Description: "List all caches",
				Execute:     inst.list,
			},
		},
	})
	return inst
}

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (c *Cache) Name() string {
	return c.tree.Node().Name
}

func (c *Cache) Description() string {
	return c.tree.Node().Description
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
	if r.Args().Len() > 1 {
		c.l.Info("clearing cache:")
		for _, value := range r.Args()[1:] {
			c.l.Info("└ " + value)
			c.cache.Get(value).Delete("")
		}
	} else {
		c.l.Info("clearing all caches")
		c.cache.Clear()
	}
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
