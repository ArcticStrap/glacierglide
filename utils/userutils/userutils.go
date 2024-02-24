package userutils

import (
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
)

func UserCan(action string, userGroups []string) bool {
	for _, v := range userGroups {
		if wikiconfig.UserGroups[v][action] {
			return true
		}
	}
	return false
}
