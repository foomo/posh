package cache

type MemoryCache map[string]MemoryNamespace

func (c MemoryCache) Clear(namespace string) {
	if namespace == "" {
		for _, value := range c {
			value.Delete("")
		}
	} else {
		c.Get(namespace).Delete("")
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
