package cache

type MemoryCache map[string]MemoryNamespace

func (c MemoryCache) Clear(namespace string) {
	if namespace == "" {
		for key := range c {
			delete(c, key)
		}
	} else {
		delete(c, namespace)
	}
}

func (c MemoryCache) Get(namespace string) Namespace {
	if _, ok := c[namespace]; !ok {
		c[namespace] = MemoryNamespace{}
	}
	return c[namespace]
}

func (c MemoryCache) List() map[string]Namespace {
	ret := map[string]Namespace{}
	for s, namespace := range c {
		ret[s] = namespace
	}
	return ret
}
