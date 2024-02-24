package wikido

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/ChaosIsFramecode/horinezumi/utils/userutils"
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
		dorouter.Post("/suspend", func(w http.ResponseWriter, r *http.Request) {
			var bReq data.SusReq

			if err := json.NewDecoder(r.Body).Decode(&bReq); err != nil {
				jsonresp.JsonERR(w, http.StatusBadRequest, "Error with decoding json: ", err)
				return
			}

			// Get editor name
			editor, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
			if err != nil {
				jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
        return
			}
			// Check if suspending a user is possible
			userGroups, err := db.GetUserGroups(editor)
			if err != nil {
			  jsonresp.JsonERR(w, http.StatusBadRequest, "Error with getting user groups: %s", err)
        return
		  }
      proceed := userutils.UserCan("suspend",userGroups)
			
      if !proceed {
				jsonresp.JsonERR(w, http.StatusUnauthorized, "Error suspending user: Permission denied", nil)
				return
			}

			err = db.SuspendUser(bReq.Target, bReq.Duration)
			if err != nil {
				jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with suspending user: %s", err)
				return
			}

			// Make response
			jsonresp.JsonOK(w, make(map[string]string), "Suspended user")
		})
		dorouter.Post("/lock", func(w http.ResponseWriter, r *http.Request) {
			// Decode lock request
			var lockReq data.LockReq

			if err := json.NewDecoder(r.Body).Decode(&lockReq); err != nil {
				jsonresp.JsonERR(w, http.StatusBadRequest, "Error with decoding json: ", err)
				return
			}

			// Redirect if not lowercase
			if strings.ToLower(lockReq.Title) != lockReq.Title {
				http.Redirect(w, r, strings.ToLower(lockReq.Title), http.StatusSeeOther)
				return
			}
			// Get editor name
			editor, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
			if err != nil {
				jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
        return
			}
			// Check if locking a page is possible
			userGroups, err := db.GetUserGroups(editor)
			if err != nil {
			  jsonresp.JsonERR(w, http.StatusBadRequest, "Error with getting user groups: %s", err)
        return
		  }
      proceed := userutils.UserCan("lock",userGroups)
			
      if !proceed {
				jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with locking page: Permission denied", nil)
				return
			}

			err = db.LockPage(lockReq.Title, int(lockReq.MinimumGroup))
			if err != nil {
				jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with locking page: %s", err)
				return
			}

			// Make response
			jsonresp.JsonOK(w, make(map[string]string), "Page Locked")
		})
	})
}
