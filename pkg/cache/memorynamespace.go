package cache

import "github.com/c-bata/go-prompt"

type MemoryNamespace map[string]interface{}

func (c MemoryNamespace) Delete(key string) {
	if key == "" {
		for key := range c {
			delete(c, key)
		}
	} else {
		delete(c, key)
	}
}

func (c MemoryNamespace) Get(key string, cb func() interface{}) interface{} {
	if _, ok := c[key]; !ok {
		if cb == nil {
			return nil
		}
		c[key] = cb()
	}
	return c[key]
}

func (c MemoryNamespace) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

func (c MemoryNamespace) GetSuggests(key string, cb func() interface{}) []prompt.Suggest {
	if v, ok := c.Get(key, cb).([]prompt.Suggest); ok {
		return v
	}
	return nil
}
