package userutils

import (
	"github.com/ChaosIsFramecode/horinezumi/utils/iputils"
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
)

func GetUserGroups(username string) []string {
	if iputils.NameIsIP(username) {
		return []string{"*"}
	}

	return []string{"*", wikiconfig.DefaultLoginGroup}
}

func UserCan(action string, userGroups []string) bool {
	for _, v := range userGroups {
		if wikiconfig.UserGroups[v][action] {
			return true
		}
	}
	return false
}
