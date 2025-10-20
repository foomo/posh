package suggests

import (
	"fmt"

	"github.com/foomo/posh/pkg/prompt/goprompt"
)

func List[T any](v []T) []goprompt.Suggest {
	ret := make([]goprompt.Suggest, len(v))
	for i, a := range v {
		ret[i] = goprompt.Suggest{Text: fmt.Sprintf("%v", a)}
	}

	return ret
}
