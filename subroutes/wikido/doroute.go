package wikido

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
  "github.com/ChaosIsFramecode/horinezumi/wikiinfo"
)

func SetupDoRoute(rt *chi.Mux, db data.Datastore) {
	// Add do subroute
	rt.Route("/d", func(dorouter chi.Router) {
		// Version info
		dorouter.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(fmt.Sprintf("Core\nHorinezumi: Version %s\nGo: Version %s\n%s: Version %s", wikiinfo.Version, runtime.Version(), db.EngineName(), db.Version())))
		})
	})
}
