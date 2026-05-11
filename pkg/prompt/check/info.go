package check

type Info struct {
	Icon   string
	Name   string
	Note   string
	Status Status
}

func NewInfo(icon, name, note string, status Status) Info {
	return Info{
		Icon:   icon,
		Name:   name,
		Note:   note,
		Status: status,
	}
}

func NewNoteInfo(icon, name, note string) Info {
	return NewInfo(icon, name, note, StatusNote)
}

func NewSuccessInfo(icon, name, note string) Info {
	return NewInfo(icon, name, note, StatusSuccess)
}

func NewWarningInfo(icon, name, note string) Info {
	return NewInfo(icon, name, note, StatusWarning)
}

func NewFailureInfo(icon, name, note string) Info {
	return NewInfo(icon, name, note, StatusFailure)
}
