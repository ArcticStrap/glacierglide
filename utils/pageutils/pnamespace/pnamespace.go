package pnamespace

import "github.com/ArcticStrap/glacierglide/wikiconfig"

const (
	Main int = iota
	MainD
	User
	UserD
	Project
	ProjectD
	File
	FileD
	Template
	TemplateD
	Help
	HelpD
	Category
	CategoryD
)

func NamespaceFromNumber(ns int) string {
	switch ns {
	case Main:
		return "main"
	case MainD:
		return "discussion"
	case User:
		return "user"
	case UserD:
		return "user_discussion"
	case Project:
		return wikiconfig.WikiName
	case ProjectD:
		return wikiconfig.WikiName + "_discussion"
	case File:
		return "file"
	case FileD:
		return "file_discussion"
	case Template:
		return "template"
	case TemplateD:
		return "template_discussion"
	case Help:
		return "help"
	case HelpD:
		return "help_discussion"
	case Category:
		return "category"
	case CategoryD:
		return "category_discussion"
	default:
		return "main"
	}
}

func NumberFromNamespace(ns string) int {
	switch ns {
	case "main":
		return Main
	case "discussion":
		return MainD
	case "user":
		return User
	case "user_discussion":
		return UserD
	case wikiconfig.WikiName:
		return Project
	case wikiconfig.WikiName + "_discussion":
		return ProjectD
	case "file":
		return File
	case "file_discussion":
		return FileD
	case "template":
		return Template
	case "template_discussion":
		return TemplateD
	case "help":
		return Help
	case "help_discussion":
		return HelpD
	case "category":
		return Category
	case "category_discussion":
		return CategoryD
	default:
		return Main
	}
}
