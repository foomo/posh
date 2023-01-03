package cache

import (
	"github.com/c-bata/go-prompt"
)

type Namespace interface {
	Get(key string, cb func() interface{}) interface{}
	Keys() []string
	Delete(key string)
	GetSuggests(key string, cb func() interface{}) []prompt.Suggest
}
