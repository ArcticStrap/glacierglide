package source

import (
	"net/http"
	"strings"

	"github.com/ArcticStrap/glacierglide/data"
)

func SetupSourceRoute(rt *http.ServeMux, db data.Datastore) {
	// Page source handler
	rt.HandleFunc("GET /api/s/{title}", func(w http.ResponseWriter, r *http.Request) {
		titleParam := r.PathValue("title")
		// Redirect if not lowercase
		if strings.ToLower(titleParam) != titleParam {
			http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
			return
		}

		// Check if page exists
		p, err := db.ReadPage(titleParam)
		if err != nil {
			w.Write([]byte("Page does not exist. Try checking your spelling if otherwise."))
			return
		}
		w.Write([]byte(p.Content))
	})
}
