package userutils

import (
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

func WhosGroupSuperior(group1 string, group2 string) string {
	index1 := 0
	index2 := 0
	curIndex := 0
	for i := range wikiconfig.UserGroups {
		if i == group1 {
			index1 = curIndex
		} else if i == group2 {
			index2 = curIndex
		}
		curIndex++
	}
	if index1 > index2 {
		return group1
	} else if index1 < index2 {
		return group2
	}
	return ""
}
