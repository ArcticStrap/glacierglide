package edit

import (
	"encoding/json"
	"net/http"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/go-chi/chi/v5"
)

func SetupEditRoute(rt *chi.Mux, db data.Datastore) {
	// Setup subrouter for wiki editing
	rt.Route("/e", func(editrouter chi.Router) {
		editrouter.Route("/{title}", func(pagerouter chi.Router) {
			// TODO: Create page
			pagerouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Handle request
				newPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&newPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with parsing json request: %s", err)
					return
				}
				// Create page in database
				if err := db.CreatePage(newPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with inserting page into database: %s", err)
					return
				}

				// Make response
				resp := make(map[string]string)
				resp["pageTitle"] = newPage.Title

				jsonresp.JsonOK(w, resp, "Page created!")
			})

			// TODO: Update page content
			pagerouter.Put("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Handle request
				uPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&uPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with parsing json request: %s", err)
					return
				}
				// Update page from database
				// Create page in database
				if err := db.UpdatePage(uPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with inserting page into database: %s", err)
					return
				}

				// Make response
				jsonresp.JsonOK(w, make(map[string]string), "Page updated!")
			})

			// TODO: Delete page
			pagerouter.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Handle request
				pageTitle := new(struct {
					Title string `json:"title"`
				})
				if err := json.NewDecoder(r.Body).Decode(&pageTitle); err != nil {
					jsonresp.JsonERR(w, 422, "Error with parsing json request: %s", err)
					return
				}
				// Delete page from database
				if err := db.DeletePage(pageTitle.Title); err != nil {
					jsonresp.JsonERR(w, 422, "Error deleting page from database: %s", err)
					return
				}

				// Make response
				jsonresp.JsonOK(w, make(map[string]string), "Page deleted!")
			})
		})
	})
}
