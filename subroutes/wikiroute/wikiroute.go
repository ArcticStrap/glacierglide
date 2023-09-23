package wikiroute

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"

	"github.com/ChaosIsFramecode/horinezumi/subroutes/editroute"
)

func SetupWikiroute(rt *chi.Mux, conn *pgx.Conn) {
	// Add view subroute
	rt.Route("/wiki", func(wikirouter chi.Router) {
		// Redirect root to main page
		wikirouter.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

		// Page view handler
		wikirouter.Route("/{title}", func(pagerouter chi.Router) {
			// Retrieve page content
			pagerouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				titleParam := chi.URLParam(r, "title")
				if titleParam == "Main_Page" {
					// Main page
					w.Write([]byte("Welcome to 堀ネズミ!"))
				} else {
					w.Write([]byte(titleParam))
				}
			})
		})
	})

	// Call edit subroute
	editroute.SetupEditRoute(rt, conn)
}
