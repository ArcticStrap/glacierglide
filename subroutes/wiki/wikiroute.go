package wiki

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupWikiroute(rt *chi.Mux, db data.Datastore) {
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
					// Redirect if not lowercase
					if strings.ToLower(titleParam) != titleParam {
						http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
						return
					}

					// Check if page exists
					p, err := db.ReadPage(titleParam)
					if err != nil {
						w.Write([]byte("Page does not exist. Try checking your spelling if otherwise."))
						return
					}

					w.Write([]byte("<h1>" + p.Title + "</h1>\r\n" + p.Content))
				}
			})
		})
	})
}
