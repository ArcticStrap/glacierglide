package handlers

import (
	"html/template"
	"net/http"

	"github.com/ArcticStrap/glacierglide/polarpages/models"
	"github.com/go-chi/chi/v5"
)

func SetupUserHandler(rt *chi.Mux) {
	// Initalize templates
	caTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/createaccount.html"))
	liTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/login.html"))
	loTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/logout.html"))
	// User routing
	rt.Get("/CreateAccount", func(w http.ResponseWriter, r *http.Request) {
		// Parse template
		caTmpl.Execute(w, struct {
			models.SessionData
			models.StaticPage
		}{UserSession(r), models.StaticPage{
			Theme: "common",
		}})
	})

	rt.Get("/Login", func(w http.ResponseWriter, r *http.Request) {
		// Parse template
		liTmpl.Execute(w, struct {
			models.SessionData
			models.StaticPage
		}{UserSession(r), models.StaticPage{
			Theme: "common",
		}})
	})

	rt.Get("/Logout", func(w http.ResponseWriter, r *http.Request) {
		http.Post("/api/Logout", "application/json", nil)
		loTmpl.Execute(w, UserSession(r))
	})
}

func UserSession(r *http.Request) models.SessionData {
	_, err := r.Cookie("gg_session")
	if err != nil {
		return models.SessionData{
			LoggedIn: false,
		}
	}

	return models.SessionData{
		LoggedIn: true,
		Username: "User",
	}
}
