package pnamespace

import "github.com/ArcticStrap/glacierglide/wikiconfig"

const (
	Main      int = 0
	MainD     int = 1
	User      int = 2
	UserD     int = 3
	Project   int = 4
	ProjectD  int = 5
	File      int = 6
	FileD     int = 7
	Template  int = 8
	TemplateD int = 9
	Help      int = 10
	HelpD     int = 11
	Category      = 12
	CategoryD     = 13
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
