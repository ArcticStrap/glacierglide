package source

import (
	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupSourceRoute(rt *chi.Mux, _ data.Datastore) {
	// Add source subroute
	rt.Route("/s", func(dorouter chi.Router) {
	})
}
