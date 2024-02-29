package do

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupDoroute(rt *chi.Mux, db data.Datastore) {
	// Add view subroute
	rt.Route("/do", func(dorouter chi.Router) {
		// Page view handler
		dorouter.Get("/version", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("Version 0.01"))
			}
    }		
	})
}
