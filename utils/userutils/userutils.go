package userutils

import (
	"github.com/ChaosIsFramecode/horinezumi/utils/iputils"
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
)

func GetUserGroups(username string) []string {
	if iputils.NameIsIP(username) {
		return []string{"*"}
	}

	return []string{"*",wikiconfig.DefaultLoginGroup}
}
