package readline

import (
	"github.com/spf13/pflag"
)

type Flags []*pflag.Flag

func (f Flags) Remove(name string) (*pflag.Flag, Flags) {
	for i, v := range f {
		if v.Name == name {
			return v, f.Splice(i, 1)
		}
	}
	return nil, f
}

func (f Flags) Slice(start, end int) Flags {
	return append(f[:start], f[end:]...)
}

func (f Flags) Splice(start, num int) Flags {
	return append(f[:start], f[start+num:]...)
}

func (f Flags) Args() Args {
	var ret Args
	for _, v := range f {
		switch v.Value.Type() {
		case "bool":
			ret = append(ret, "--"+v.Name)
		default:
			ret = append(ret, "--"+v.Name, v.Value.String())
		}
	}
	return ret
}
