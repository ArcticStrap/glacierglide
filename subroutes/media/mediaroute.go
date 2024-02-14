package media

import (
	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupMediaRoute(rt *chi.Mux, _ data.Datastore) {
	// Add media subroute
	rt.Route("/m", func(dorouter chi.Router) {
	})
}
