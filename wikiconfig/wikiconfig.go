package wikiconfig

// Main settings

const WikiName string = "arcticstrap"

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
		"lock":    true,
	},
	// TODO: "techie":        {},
	"administrator": {
		"managegroup": true,
		"renameuser":  true,
	},
	"manager": {},
}

var GroupAssignRights map[string][]string = map[string][]string{
	"moderator":     {},
	"administrator": {"moderator"},
	"manager":       {"administrator", "moderator"},
}
