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
	lTmpl := template.Must(template.ParseFiles("polarpages/templates/base.html", "polarpages/templates/login.html"))
	// User routing
	rt.Get("/CreateAccount", func(w http.ResponseWriter, _ *http.Request) {
		// Parse template
		caTmpl.Execute(w, models.StaticPage{
			Theme: "common",
		})
	})

	rt.Get("/Login", func(w http.ResponseWriter, _ *http.Request) {
		// Parse template
		lTmpl.Execute(w, models.StaticPage{
			Theme: "common",
		})
	})
}
