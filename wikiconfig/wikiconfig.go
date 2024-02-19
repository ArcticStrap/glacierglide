package wikiconfig

var UserGroups map[string]map[string]bool = map[string]map[string]bool{
	"*": {
		"createaccount": true,
		"delete":        true,
	},
}
