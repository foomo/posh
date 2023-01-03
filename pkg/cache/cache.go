package cache

type Cache interface {
	Get(namespace string) Namespace
	Clear(namespace string)
	List() map[string]Namespace
}
