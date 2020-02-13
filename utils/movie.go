package utils

import (
	"github.com/greenac/artemis/logger"
	"strings"
)

func AddFollowingUnderscore(an string, mn string) string {
	i := strings.Index(mn, an)
	if i == -1 {
		return mn
	}

	mnr := []rune(mn)
	ti := i + len(an)
	if ti > len(mnr)-1 {
		return mn
	}

	if mnr[ti] == '_' || mnr[ti] == '.' {
		return mn
	}

	return addUnderscoreAtIndex(&mnr, ti)
}

func AddPrecedingUnderscore(an string, mn string) string {
	i := strings.Index(mn, an)
	if i == -1 {
		return mn
	}

	mnr := []rune(mn)
	if i == 0 {
		return mn
	}

	logger.Log("character at:", i, "is:", string(mnr[i]))

	if mnr[i-1] == '_' {
		return mn
	}

	return addUnderscoreAtIndex(&mnr, i)
}

func AddMiddleUnderscore(n1 string, n2 string, mn string) string {
	i1 := strings.Index(mn, n1)
	if i1 == -1 {
		return mn
	}

	i2 := strings.Index(mn, n1)
	if i2 == -1 {
		return mn
	}

	mnr := []rune(mn)
	if len(n1) + len(n2) >= len(mnr) {
		return mn
	}

	if i1+len(n1) != i2 {
		return mn
	}


	return addUnderscoreAtIndex(&mnr, i2)
}

func addUnderscoreAtIndex(rp *[]rune, i int) string {
	mnr := *rp
	mnr = append(mnr[:i], append([]rune{'_'}, mnr[i:]...)...)

	return string(mnr)
}