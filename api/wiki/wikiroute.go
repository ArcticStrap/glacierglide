package wiki

import (
	"net/http"
	"strings"

	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/mprender"
	"github.com/ArcticStrap/glacierglide/mprender/markdown"
)

func SetupWikiRoute(rt *http.ServeMux, db data.Datastore) {
	// Add view subroute
		// Redirect root to main page
		rt.Handle("GET /", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther))

		// Page view handler
		rt.HandleFunc("GET /{title}", func(w http.ResponseWriter, r *http.Request) {
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
