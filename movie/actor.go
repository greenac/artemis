package movie

import "github.com/greenac/artemis/tools"

type Actor struct {
	FirstName *string
	LastName *string
	MiddleName *string
	Files *map[string]tools.File
}

func (act *Actor)GetFirstName() string {
	if act.FirstName == nil {
		return ""
	}

	return *act.FirstName
}

func (act *Actor)GetMiddleName() string {
	if act.MiddleName == nil {
		return ""
	}

	return *act.MiddleName
}

func (act *Actor)GetLastName() string {
	if act.LastName == nil {
		return ""
	}

	return *act.LastName
}

func (act *Actor)FullName() string {
	name := ""
	if act.FirstName != nil {
		name += *act.FirstName
	}

	if act.MiddleName != nil {
		if name == "" {
			name = *act.MiddleName
		} else {
			name += "_" + *act.MiddleName
		}
	}

	if act.LastName != nil {
		if name == "" {
			name = *act.LastName
		} else {
			name += "_" + *act.LastName
		}
	}

	return name
}

func (act *Actor)addFile(f *tools.File) {
	if act.Files == nil {
		files := make(map[string]tools.File)
		files[*f.Name()] = *f
		act.Files = &files
	} else {
		(*act.Files)[*f.Name()] = *f
	}
}

func (act *Actor)AddFiles(fls *map[string]tools.File) {
	if fls == nil {
		return
	}

	for _, f := range *fls {
		act.addFile(&f)
	}
}
