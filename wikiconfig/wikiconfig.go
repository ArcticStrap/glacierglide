package wikiconfig

// User related settings

const DefaultLoginGroup string = "user"

var UserGroups map[string]map[string]bool = map[string]map[string]bool{
	"*": {
		"createaccount": true,
		"edit":          true,
		"delete":        true,
		"create":        true,
	},
	"user": {},
	// TODO: "affirmed":      {},
	// TODO: "trusted":       {},
	// TODO: "bot":           {},
	"moderator": {
		"suspend": true,
	},
	// TODO: "techie":        {},
	"administrator": {
		"managegroup": true,
		"renameuser":  true,
	},
	"manager": {},
}
