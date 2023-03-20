package readline

import (
	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
}

func NewFlagSet(name string) *FlagSet {
	fs := pflag.NewFlagSet(name, pflag.ContinueOnError)
	fs.ParseErrorsWhitelist = pflag.ParseErrorsWhitelist{UnknownFlags: true}
	return &FlagSet{
		FlagSet: fs,
	}
}

func (s *FlagSet) Visited() Flags {
	var ret Flags
	s.Visit(func(f *pflag.Flag) {
		ret = append(ret, f)
	})
	return ret
}

func (s *FlagSet) SetValues(name string, values ...string) error {
	return s.SetAnnotation(name, "values", values)
}

func (s *FlagSet) GetValues(name string) []string {
	if f := s.FlagSet.Lookup(name); f == nil {
		return nil
	} else if v, ok := f.Annotations["values"]; ok {
		return v
	}
	return nil
}
