package wikido

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupDoroute(rt *chi.Mux, _ data.Datastore) {
	// Add view subroute
	rt.Route("/d", func(dorouter chi.Router) {
		// Page view handler
		dorouter.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("Version 0.01"))
		})
	})
}
