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
