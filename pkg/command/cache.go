package command

import (
	"context"
	"fmt"
	"sort"

	"github.com/foomo/posh/pkg/agent"
	"github.com/foomo/posh/pkg/cache"
	"github.com/foomo/posh/pkg/command/tree"
	"github.com/foomo/posh/pkg/log"
	"github.com/foomo/posh/pkg/prompt/goprompt"
	"github.com/foomo/posh/pkg/readline"
	"github.com/foomo/posh/pkg/util/suggests"
	"github.com/pterm/pterm"
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
			{
				Name:        "get",
				Description: "Get a cached value",
				Args: tree.Args{
					{
						Name:        "Namespace",
						Description: "Name of the namespace.",
						Suggest: func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
							return suggests.List(lo.Keys(inst.cache.List()))
						},
					},
					{
						Name:        "Key",
						Description: "Name of the cached key.",
						Suggest: func(ctx context.Context, t tree.Root, r *readline.Readline) []goprompt.Suggest {
							return suggests.List(inst.cache.Get(r.Args().At(1)).Keys())
						},
					},
				},
				Execute: inst.get,
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

// Describe implements the optional Describer interface, letting
// `posh agent catalog` describe this command's subtree.
func (c *Cache) Describe(ctx context.Context) CommandInfo {
	return c.tree.Describe(ctx)
}

// Skill implements the optional Skiller interface. The catalog already carries
// the structure; what it cannot show is that the namespaces are populated at
// runtime, which is what makes `cache list` the necessary first step.
func (c *Cache) Skill(ctx context.Context, name string) string {
	return "#### Notes\n\n" +
		"The cache is in-memory and per shell process: it is empty on start and is\n" +
		"lost on exit, so it never needs clearing between sessions.\n\n" +
		"Namespaces are created by whichever command populates them, so they cannot\n" +
		"be listed ahead of time - run `cache list` to discover the namespaces and\n" +
		"keys that currently exist, then `cache get [Namespace] [Key]` to read one\n" +
		"value. `cache list` prints keys only; it never prints values.\n\n" +
		"`cache clear` with no argument drops every namespace. Prefer naming the\n" +
		"namespaces to clear - it accepts several - so unrelated cached lookups are\n" +
		"not discarded with them. Clearing is safe: the next command that needs a\n" +
		"value recomputes it, the only cost is the recomputation."
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

// get prints a single cached value. Values are otherwise only visible in the
// list at trace level, so this is the way to inspect one.
func (c *Cache) get(ctx context.Context, r *readline.Readline) error {
	namespace, key := r.Args().At(1), r.Args().At(2)
	// A nil callback returns the cached value without populating it.
	value := c.cache.Get(namespace).Get(key, nil)

	return agent.Render(
		func() any { return CacheGet{Namespace: namespace, Key: key, Value: value} },
		func() error {
			c.l.Info(fmt.Sprintf("%v", value))
			return nil
		},
	)
}

func (c *Cache) list(ctx context.Context, r *readline.Readline) error {
	cacheList := c.cache.List()
	namespaces := lo.Keys(cacheList)
	sort.Strings(namespaces)

	return agent.Render(
		func() any { return c.listPayload(cacheList, namespaces) },
		func() error { return c.listTree(cacheList, namespaces) },
	)
}

// listPayload builds the JSON form. Only keys are included, mirroring the tree
// at default level - use `cache get` to read a value.
//
// It is namespace-keyed rather than reusing agent.Tree: the generic
// {"text","children"} shape would make a consumer index into anonymous nodes to
// find a namespace.
func (c *Cache) listPayload(cacheList map[string]cache.Namespace, namespaces []string) CacheList {
	values := make([]CacheNamespace, 0, len(namespaces))

	for _, ns := range namespaces {
		keys := cacheList[ns].Keys()
		sort.Strings(keys)

		if keys == nil {
			keys = []string{}
		}

		values = append(values, CacheNamespace{Namespace: ns, Keys: keys})
	}

	return CacheList{Namespaces: values}
}

// listTree renders the human-formatted tree, including cached values at trace
// level.
func (c *Cache) listTree(cacheList map[string]cache.Namespace, namespaces []string) error {
	list := pterm.LeveledList{}

	for _, ns := range namespaces {
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

	return agent.Tree(list)
}
