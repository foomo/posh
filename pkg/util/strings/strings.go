package strings

import (
	"math"
	"strings"
)

func PadEnd(s, suffix string, count int) string {
	return s + strings.Repeat(suffix, int(math.Max(0, float64(count-len(s)))))
}
