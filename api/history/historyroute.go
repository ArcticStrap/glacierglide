package history

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/jsonresp"
)

func SetupHistoryRoute(rt *http.ServeMux, db data.Datastore) {
	// History subroute handler
		// Retrieve page history
		rt.HandleFunc("GET /h/{title}", func(w http.ResponseWriter, r *http.Request) {
			// Expect json response
			w.Header().Set("Content-Type", "application/json")

			titleParam := r.PathValue("title")
			// Redirect if not lowercase
			if strings.ToLower(titleParam) != titleParam {
				http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
				return
			}

			// Fetch page history
			pH, err := db.FetchPageHistory(titleParam)
			if err != nil {
				jsonresp.JsonERR(w, http.StatusBadRequest, "Error with fetching page history: %s", err)
				return
			}

			// Convert list to json
			json.NewEncoder(w).Encode(pH)
		})
}
