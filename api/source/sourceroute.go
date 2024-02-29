package source

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupSourceRoute(rt chi.Router, db data.Datastore) {
	// Add source subroute
	rt.Route("/s", func(sourcerouter chi.Router) {
		// Page source handler
		sourcerouter.Get("/{title}", func(w http.ResponseWriter, r *http.Request) {
			titleParam := chi.URLParam(r, "title")
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
			w.Write([]byte(p.Content))
		})
	})
}
