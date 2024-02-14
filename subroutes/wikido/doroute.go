package wikido

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
)

func SetupDoroute(rt *chi.Mux, db data.Datastore) {
	// Add view subroute
	rt.Route("/d", func(dorouter chi.Router) {
		// Page view handler
		dorouter.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
      w.Write([]byte(fmt.Sprintf("Core\nHorinezumi: Version 0.01\nGo: Version %s\n%s: Version %s", runtime.Version(),db.EngineName(),db.Version())))
		})
	})
}
