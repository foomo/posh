package readline

type Args []string

// ------------------------------------------------------------------------------------------------
// ~ Public methods
// ------------------------------------------------------------------------------------------------

func (a Args) At(v int) string {
	if a.HasIndex(v) {
		return a[v]
	} else {
		return ""
	}
}

func (a Args) AtDefault(v int, fallback string) string {
	if a.HasIndex(v) {
		return a[v]
	} else {
		return fallback
	}
}

func (a Args) Shift() (string, Args) {
	if a.HasIndex(0) {
		return a[0], a[1:]
	} else {
		return "", nil
	}
}

func (a Args) Empty() bool {
	return a == nil || a.LenIs(0)
}

func (a Args) Has(v string) bool {
	for _, arg := range a {
		if arg == v {
			return true
		}
	}
	return false
}

func (a Args) HasIndex(v int) bool {
	return a.LenGte(v + 1)
}

func (a Args) Len() int {
	return len(a)
}

func (a Args) LenIs(v int) bool {
	return a.Len() == v
}

func (a Args) LenGt(v int) bool {
	return a.Len() > v
}

func (a Args) LenGte(v int) bool {
	return a.Len() >= v
}

func (a Args) LenLt(v int) bool {
	return a.Len() < v
}

func (a Args) LenLte(v int) bool {
	return a.Len() <= v
}

func (a Args) Last() string {
	if !a.Empty() {
		return a[a.Len()-1]
	} else {
		return ""
	}
}

func (a Args) IndexOf(v string) int {
	for i, s := range a {
		if s == v {
			return i
		}
	}
	return -1
}

func (a Args) From(start int) Args {
	return a[start:]
}

func (a Args) To(end int) Args {
	return a[:end]
}

func (a Args) Slice(start, end int) Args {
	return append(a[:start], a[end:]...)
}

func (a Args) Splice(start, num int) Args {
	return append(a[:start], a[start+num:]...)
}
