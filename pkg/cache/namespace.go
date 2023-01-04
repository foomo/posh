package cache

import (
	"github.com/c-bata/go-prompt"
)

type Namespace interface {
	Get(key string, cb func() any) any
	Keys() []string
	Delete(key string)
	GetSuggests(key string, cb func() any) []prompt.Suggest
}
