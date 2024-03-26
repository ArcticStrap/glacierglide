package wiki

import (
	"net/http"
	"strings"

	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/mprender"
	"github.com/ArcticStrap/glacierglide/mprender/markdown"
)

func SetupWikiRoute(rt *http.ServeMux, db data.Datastore) {
	// Page view handler
	rt.HandleFunc("GET /api/wiki/{title}", func(w http.ResponseWriter, r *http.Request) {
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

		// Check page format
		switch p.MPType {
		case mprender.Markdown:
			p.Content = markdown.ToHTML(p.Content)
		}

		w.Write([]byte(p.Content))
	})
}
