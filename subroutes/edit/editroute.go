package edit

import (
	"encoding/json"
	"net/http"
	"strings"

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

				// Fetch title param
				titleParam := chi.URLParam(r, "title")
				// Redirect if not lowercase
				if strings.ToLower(titleParam) != titleParam {
					http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
					return
				}

				// Handle request
				newPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&newPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with parsing json request: %s", err)
					return
				}

				// Add title if nil
				if newPage.Title == "" {
					newPage.Title = titleParam
				}
				// Lowercase page title
				newPage.Title = strings.ToLower(newPage.Title)

				// Check if page already exists
				pageExists, _ := db.ReadPage(newPage.Title)
				if pageExists != nil {
					jsonresp.JsonOK(w, make(map[string]string), "Page already exists!")
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

			// Update page content
			pagerouter.Put("/", data.CallJWTAuth(db, false, func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Fetch title param
				titleParam := chi.URLParam(r, "title")
				// Redirect if not lowercase
				if strings.ToLower(titleParam) != titleParam {
					http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
					return
				}

				// Handle request
				uPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&uPage); err != nil {
					jsonresp.JsonERR(w, 422, "Error with parsing json request: %s", err)
					return
				}
				// Add title if nil
				if uPage.Title == "" {
					uPage.Title = titleParam
				}
				// Lowercase page title
				uPage.Title = strings.ToLower(uPage.Title)

				// Authenticate token, default to ip address
				tokenStr := r.Header.Get("authtoken")
				editor := "0.0.0.0"
				if tokenStr != "" {
					token, err := data.ValidateJWT(tokenStr)
					if err != nil || !token.Valid {
						jsonresp.JsonERR(w, 401, "Invalid token", nil)
						return
					}

					claims, ok := token.Claims.(*data.UserClaims)
					if !ok {
						jsonresp.JsonERR(w, 401, "Invalid token claims", nil)
						return
					}

					editor, err = db.GetUsernameFromId(claims.UserID)
					if err != nil {
						jsonresp.JsonERR(w, 401, "Failed to get user account: %s", err)
						return
					}
				} else {
					// Get ip address if not logged in
					editor = strings.Split(r.RemoteAddr, ":")[0]
				}

				// Update page from database
				if err := db.UpdatePage(uPage, editor); err != nil {
					jsonresp.JsonERR(w, 422, "Error with inserting page into database: %s", err)
					return
				}

				// Make response
				jsonresp.JsonOK(w, make(map[string]string), "Page updated!")
			}))

			// TODO: Delete page
			pagerouter.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				// Expect json response
				w.Header().Set("Content-Type", "application/json")

				// Fetch title param
				titleParam := chi.URLParam(r, "title")
				// Redirect if not lowercase
				if strings.ToLower(titleParam) != titleParam {
					http.Redirect(w, r, strings.ToLower(titleParam), http.StatusSeeOther)
					return
				}

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
