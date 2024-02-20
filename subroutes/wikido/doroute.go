package wikido

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/ChaosIsFramecode/horinezumi/utils/userutils"
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
	"github.com/ChaosIsFramecode/horinezumi/wikiinfo"
)

func SetupDoRoute(rt *chi.Mux, db data.Datastore) {
	// Add do subroute
	rt.Route("/d", func(dorouter chi.Router) {
    // GENERAL REQUESTS
		// Version info
		dorouter.Get("/version", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(fmt.Sprintf("Core\nHorinezumi: Version %s\nGo: Version %s\n%s: Version %s", wikiinfo.Version, runtime.Version(), db.EngineName(), db.Version())))
		})
    // MODERATION REQUESTS
    // TODO
    dorouter.Post("/suspend",func(w http.ResponseWriter, r *http.Request) {

    })
    dorouter.Post("/lock",func(w http.ResponseWriter, r *http.Request) {
        // Fetch title param
				titleParam := chi.URLParam(r, "title")
				// Redirect if not lowercase
				if strings.ToLower(titleParam) != titleParam {
					http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
					return
				}
				// Get editor name
				_, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
				}
				// Check if locking a page is possible
				userGroups := userutils.GetUserGroups(r.RemoteAddr)
				proceed := false
				for _, v := range userGroups {
					if wikiconfig.UserGroups[v]["delete"] {
						proceed = true
						break
					}
				}
				if !proceed {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with locking page: Permission denied", nil)
					return
				}

    })
	})
}
