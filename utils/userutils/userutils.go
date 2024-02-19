package userutils

import "github.com/ChaosIsFramecode/horinezumi/utils/iputils"

func GetUserGroups(username string) []string {
	if iputils.NameIsIP(username) {
		return []string{"*"}
	}

	return []string{"*"}
}
