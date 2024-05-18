package userutils

import (
	"slices"

	"github.com/ArcticStrap/glacierglide/wikiconfig"
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

func ValidRightsReq(userGroup string, add []string, remove []string) bool {
	for _, v := range add {
		if !slices.Contains(wikiconfig.GroupAssignRights[userGroup], v) {
			return false
		}
	}
	for _, v := range remove {
		if !slices.Contains(wikiconfig.GroupAssignRights[userGroup], v) {
			return false
		}
	}
	return true
}
