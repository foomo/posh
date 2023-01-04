package suggests

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func List[T any](v []T) []prompt.Suggest {
	ret := make([]prompt.Suggest, len(v))
	for i, a := range v {
		ret[i] = prompt.Suggest{Text: fmt.Sprintf("%v", a)}
	}
	return ret
}
