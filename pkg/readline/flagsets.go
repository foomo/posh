package readline

import (
	"github.com/spf13/pflag"
)

type FlagSets struct {
	sets map[string]*FlagSet
}

func NewFlagSets() *FlagSets {
	return &FlagSets{
		sets: map[string]*FlagSet{},
	}
}

func (s *FlagSets) Default() *FlagSet {
	return s.Get("default")
}

func (s *FlagSets) Internal() *FlagSet {
	return s.Get("internal")
}

func (s *FlagSets) Get(name string) *FlagSet {
	if _, ok := s.sets[name]; !ok {
		s.sets[name] = NewFlagSet(name)
	}

	return s.sets[name]
}

func (s *FlagSets) Parse(arguments []string) error {
	for _, set := range s.sets {
		if err := set.Parse(arguments); err != nil {
			return err
		}
	}

	return nil
}

func (s *FlagSets) All() *FlagSet {
	fs := NewFlagSet("all")
	for _, set := range s.sets {
		fs.AddFlagSet(set.FlagSet)
	}

	return fs
}

func (s *FlagSets) Visit(fn func(*pflag.Flag)) Flags {
	var ret Flags

	for _, set := range s.sets {
		set.Visit(fn)
	}

	return ret
}

func (s *FlagSets) VisitAll(fn func(*pflag.Flag)) Flags {
	var ret Flags

	for _, set := range s.sets {
		set.VisitAll(fn)
	}

	return ret
}

func (s *FlagSets) Visited() Flags {
	var ret Flags
	for _, set := range s.sets {
		ret = append(ret, set.Visited()...)
	}

	return ret
}

func (s *FlagSets) ParseAll(arguments []string, fn func(flag *pflag.Flag, value string) error) error {
	for _, group := range s.sets {
		if err := group.ParseAll(arguments, fn); err != nil {
			return err
		}
	}

	return nil
}
