package tree

type Args []*Arg

func (a Args) Last() *Arg {
	if len(a) > 0 {
		return a[len(a)-1]
	} else {
		return nil
	}
}
