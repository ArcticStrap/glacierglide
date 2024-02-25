package userutils

import (
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
)

func GroupFromIndex(id int) string {
	v := 0
	for i := range wikiconfig.UserGroups {
		if v == id {
			return i
		}
		v++
	}
	return ""
}

func UserCan(action string, userGroups []string) bool {
	for _, v := range userGroups {
		if wikiconfig.UserGroups[v][action] {
			return true
		}
	}
	return false
}
