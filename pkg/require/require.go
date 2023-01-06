package require

import (
	"github.com/foomo/fender/fend"
	"github.com/foomo/posh/pkg/log"
)

func First(l log.Logger, fends ...any) error {
	var allFends []fend.Fend
	for _, value := range fends {
		switch v := value.(type) {
		case fend.Fend:
			allFends = append(allFends, v)
		case []fend.Fend:
			allFends = append(allFends, v...)
		default:
			l.Warn("unknown type:", v)
		}
	}
	if fendErr, err := fend.First(allFends...); err != nil {
		return err
	} else if fendErr != nil {
		return fendErr
	}
	return nil
}
