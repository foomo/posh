package cache

import (
	"sync"
)

type MemoryCache struct {
	store sync.Map
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		store: sync.Map{},
	}
}

func (c *MemoryCache) Clear(namespaces ...string) {
	if len(namespaces) == 0 {
		c.store.Range(func(key, value interface{}) bool {
			namespaces = append(namespaces, key.(string)) //nolint:forcetypeassert
			return true
		})
	}
	for _, namespace := range namespaces {
		c.Get(namespace).Delete()
	}
}

func (c *MemoryCache) Get(namespace string) Namespace {
	value, _ := c.store.LoadOrStore(namespace, &MemoryNamespace{
		store: sync.Map{},
	})
	return value.(*MemoryNamespace) //nolint:forcetypeassert
}

func (c *MemoryCache) List() map[string]Namespace {
	ret := map[string]Namespace{}
	c.store.Range(func(k, v interface{}) bool {
		ret[k.(string)] = v.(*MemoryNamespace) //nolint:forcetypeassert
		return true
	})
	return ret
}
