package media

import (
	"github.com/go-chi/chi/v5"

	"github.com/ArcticStrap/glacierglide/data"
)

func SetupMediaRoute(rt chi.Router, _ data.Datastore) {
	// Add media subroute
	rt.Route("/m", func(dorouter chi.Router) {
	})
}
