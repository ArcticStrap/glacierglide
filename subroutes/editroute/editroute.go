package editroute

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupEditRoute(rt *chi.Mux) {
	// Setup subrouter for wiki editing
	rt.Route("/e", func(sr chi.Router) {
		sr.Route("/{title}", func(pr chi.Router) {
			// TODO: Create page
			pr.Post("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Make response
				resp := make(map[string]string)
				resp["message"] = "Page created"
				jsonResp, err := json.Marshal(resp)
				if err != nil {
					log.Fatalf("Error happened in JSON marshal. Err: %s", err)
				}

				w.Write(jsonResp)
			})

			// TODO: Update page content
			pr.Put("/", func(w http.ResponseWriter, r *http.Request) {
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
			pr.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

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
