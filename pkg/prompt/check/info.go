package check

type Info struct {
	Name   string
	Note   string
	Status Status
}

func NewNoteInfo(name, note string) Info {
	return Info{
		Name:   name,
		Note:   note,
		Status: StatusNote,
	}
}

func NewSuccessInfo(name, note string) Info {
	return Info{
		Name:   name,
		Note:   note,
		Status: StatusSuccess,
	}
}

func NewFailureInfo(name, note string) Info {
	return Info{
		Name:   name,
		Note:   note,
		Status: StatusFailure,
	}
}
