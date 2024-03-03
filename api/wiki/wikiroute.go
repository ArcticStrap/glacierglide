package wiki

import (
	"net/http"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-chi/chi/v5"

	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/mprender"
	"github.com/ArcticStrap/glacierglide/mprender/markdown"
)

func SetupWikiRoute(rt chi.Router, db data.Datastore) {
	// Add view subroute
	rt.Route("/wiki", func(wikirouter chi.Router) {
		// Redirect root to main page
		wikirouter.Get("/", http.RedirectHandler("/wiki/Main_Page", http.StatusSeeOther).ServeHTTP)

		// Page view handler
		wikirouter.Get("/{title}", func(w http.ResponseWriter, r *http.Request) {
			titleParam := chi.URLParam(r, "title")
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

      // Capitalize first letter
      sr, size := utf8.DecodeRuneInString(p.Title)
      if sr != utf8.RuneError {
        p.Title = string(unicode.ToUpper(sr)) + p.Title[size:]
      }


			// Check page format
			switch p.MPType {
			case mprender.Markdown:
				p.Content = markdown.ToHTML(p.Content)
			}

			w.Write([]byte("<h1>" + p.Title + "</h1><br/>" + p.Content))
		})
	})
}
