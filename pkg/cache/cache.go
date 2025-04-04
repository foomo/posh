package cache

type Cache interface {
	Get(namespace string) Namespace
	Clear(namespaces ...string)
	List() map[string]Namespace
}
