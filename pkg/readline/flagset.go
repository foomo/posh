package readline

import (
	"strconv"

	"github.com/spf13/pflag"
)

type FlagSet struct {
	*pflag.FlagSet
}

func NewFlagSet() *FlagSet {
	return &FlagSet{
		FlagSet: pflag.NewFlagSet("readline", pflag.ContinueOnError),
	}
}

func (a *FlagSet) SetValues(name string, values ...string) error {
	return a.SetAnnotation(name, "values", values)
}

func (a *FlagSet) GetValues(name string) []string {
	if f := a.FlagSet.Lookup(name); f == nil {
		return nil
	} else if v, ok := f.Annotations["values"]; ok {
		return v
	}
	return nil
}

func (a *FlagSet) GetString(name string) string {
	if f := a.FlagSet.Lookup(name); f == nil {
		return ""
	} else if !a.flagIsSet(name) {
		return f.DefValue
	} else {
		return f.Value.String()
	}
}

func (a *FlagSet) GetInt64(name string) int64 {
	if value := a.GetString(name); value == "" {
		return 0
	} else if v, err := strconv.ParseInt(value, 10, 64); err != nil {
		return 0
	} else {
		return v
	}
}

func (a *FlagSet) GetFloat64(name string) float64 {
	if value := a.GetString(name); value == "" {
		return 0
	} else if v, err := strconv.ParseFloat(value, 64); err != nil {
		return 0
	} else {
		return v
	}
}

func (a *FlagSet) GetBool(name string) bool {
	if value := a.GetString(name); value == "" {
		return false
	} else if v, err := strconv.ParseBool(value); err != nil {
		return false
	} else {
		return v
	}
}

func (a *FlagSet) flagIsSet(name string) bool {
	found := false
	if fs := a.FlagSet; fs != nil {
		fs.Visit(func(f *pflag.Flag) {
			if f.Name == name {
				found = true
			}
		})
	}
	return found
}
