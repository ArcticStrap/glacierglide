package editroute

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/go-chi/chi/v5"
)

func SetupEditRoute(rt *chi.Mux, db *data.PostgresBase) {
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
					log.Fatalf("Error with parsing json request: %s", err)
					return
				}
				// Create page in database
				if err := db.CreatePage(newPage); err != nil {
					log.Fatalf("Error with inserting page into database: %s", err)
					return
				}

				// Make response
				resp := make(map[string]string)
				resp["message"] = "Page created"
				resp["pageTitle"] = newPage.Title
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}

				w.Write(jsonResp)
			})

			// TODO: Update page content
			pagerouter.Put("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Make response
				resp := make(map[string]string)
				resp["message"] = "Page updated"
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}

				w.Write(jsonResp)
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
					log.Fatalf("Error with parsing json request: %s", err)
					return
				}
				// Delete page from database
				if err := db.DeletePage(pageTitle.Title); err != nil {
					log.Fatalf("Error deleting page from database: %s", err)
					return
				}

				// Make response
				resp := make(map[string]string)
				resp["message"] = "Page deleted"
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}

				w.Write(jsonResp)
			})
		})
	})
}
