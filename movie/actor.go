package movie

type Actor struct {
	FirstName *string
	LastName *string
	MiddleName *string
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
