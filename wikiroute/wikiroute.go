package wikiroute

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupWikiroute(rt *chi.Mux) {
	// Add subroute
	rt.Route("/wiki", func(sr chi.Router) {
		// Redirect root to main page
		sr.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

		// Page handler
		sr.Get("/{title}", func(w http.ResponseWriter, r *http.Request) {
			titleParam := chi.URLParam(r, "title")
			if titleParam == "Main_Page" {
				// Main page
				w.Write([]byte("Welcome to 堀ネズミ!"))
			} else {
				w.Write([]byte(titleParam))
			}
		})
	})
}
