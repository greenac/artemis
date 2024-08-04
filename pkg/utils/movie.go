package utils

import (
	"github.com/greenac/artemis/pkg/artemiserror"
	"strings"
)

func AddFollowingUnderscore(an string, mn string) (string, error) {
	i := strings.Index(mn, an)
	if i == -1 {
		return mn, artemiserror.New(artemiserror.InvalidName)
	}

	mnr := []rune(mn)
	ti := i + len(an)
	if ti > len(mnr)-1 {
		return mn, nil
	}

	if mnr[ti] == '_' || mnr[ti] == '.' {
		return mn, nil
	}

	return addUnderscoreAtIndex(&mnr, ti), nil
}

func AddPrecedingUnderscore(an string, mn string) (string, error) {
	i := strings.Index(mn, an)
	if i == -1 {
		return mn, artemiserror.New(artemiserror.InvalidName)
	}

	mnr := []rune(mn)
	if i == 0 {
		return mn, nil
	}

	if mnr[i-1] == '_' {
		return mn, nil
	}

	return addUnderscoreAtIndex(&mnr, i), nil
}

func AddMiddleUnderscore(n1 string, n2 string, mn string) (string, error) {
	i1 := strings.Index(mn, n1)
	if i1 == -1 {
		return mn, artemiserror.New(artemiserror.InvalidName)
	}

	i2 := strings.Index(mn, n1)
	if i2 == -1 {
		return mn, artemiserror.New(artemiserror.InvalidName)
	}

	mnr := []rune(mn)
	if len(n1)+len(n2) >= len(mnr) {
		return mn, nil
	}

	if i1+len(n1) != i2 {
		return mn, nil
	}

	return addUnderscoreAtIndex(&mnr, i2), nil
}

func addUnderscoreAtIndex(rp *[]rune, i int) string {
	mnr := *rp
	mnr = append(mnr[:i], append([]rune{'_'}, mnr[i:]...)...)

	return string(mnr)
}

func AddTailingNameToMovie(movieName string, actorName string) (string, error) {
	pts := strings.Split(movieName, ".")
	if len(pts) != 2 {
		return movieName, artemiserror.New(artemiserror.InvalidName)
	}

	mp := []rune(pts[0])

	if mp[len(mp)-1] != '_' {
		mp = append(mp, '_')
	}

	return string(mp) + actorName + "." + pts[1], nil
}

func IsNameFormatCorrect(movieName string, actorName string) bool {
	index := strings.Index(movieName, actorName)
	if index == -1 {
		return false
	}

	mp := []rune(movieName)

	ft := (index == 0 || mp[index-1] == '_') &&
		(mp[index+len(actorName)] == '.' || mp[index+len(actorName)] == '_')
	return ft
}

func AddNameToMovieAfterName(mn string, name string, targetName string) (string, error) {
	i := strings.Index(mn, targetName)
	if i == -1 {
		return "", artemiserror.New(artemiserror.InvalidName)
	}

	rns := []rune(mn)
	t := i + len(targetName)
	if rns[t] != '_' {
		rns = append(rns[:t], append([]rune{'_'}, rns[t:]...)...)
	}

	t += 1

	j := strings.Index(string(rns), name)
	if j == -1 || j != t {
		rns = append(rns[:t], append([]rune(name), rns[t:]...)...)
	}

	return AddFollowingUnderscore(name, string(rns))
}
