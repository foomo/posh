package cache

import (
	"sync"

	"github.com/c-bata/go-prompt"
)

type MemoryNamespace struct {
	store sync.Map
}

func (c *MemoryNamespace) Delete(keys ...string) {
	if len(keys) == 0 {
		c.store.Clear()
	} else {
		for _, key := range keys {
			c.store.Delete(key)
		}
	}
}

func (c *MemoryNamespace) Get(key string, cb func() any) any {
	value, ok := c.store.Load(key)
	if !ok && cb != nil {
		value = cb()
		c.store.Store(key, value)
	}
	return value
}

func (c *MemoryNamespace) Keys() []string {
	var keys []string
	c.store.Range(func(k, v interface{}) bool {
		keys = append(keys, k.(string))
		return true
	})
	return keys
}

func (c *MemoryNamespace) GetSuggests(key string, cb func() any) []prompt.Suggest {
	if v, ok := c.Get(key, cb).([]prompt.Suggest); ok {
		return v
	}
	return nil
}
