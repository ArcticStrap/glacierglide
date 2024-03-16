package namespace

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
