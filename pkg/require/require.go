package require

import (
	"context"

	"github.com/foomo/fender"
	"github.com/foomo/fender/fend"
	"github.com/foomo/posh/pkg/log"
)

func First(ctx context.Context, l log.Logger, fends ...interface{}) error {
	var allFends fend.Fends
	for _, value := range fends {
		switch v := value.(type) {
		case fend.Fend:
			allFends = allFends.Append(v)
		case fend.Fends:
			allFends = allFends.Append(v...)
		case []fend.Fend:
			allFends = allFends.Append(v...)
		default:
			l.Warn("unknown type", v)
		}
	}
	if err := fender.First(ctx, allFends...); err != nil {
		return err
	}
	return nil
}
