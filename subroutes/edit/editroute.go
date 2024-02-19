package edit

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ChaosIsFramecode/horinezumi/appsignals"
	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/ChaosIsFramecode/horinezumi/utils/userutils"
	"github.com/ChaosIsFramecode/horinezumi/wikiconfig"
	"github.com/go-chi/chi/v5"
)

// Authenticate token, default to ip address if empty token.
func GetLoginStatus(tokenStr string, r *http.Request, db data.Datastore) (string, error) {
	if tokenStr != "" {
		token, err := data.ValidateJWT(tokenStr)
		if err != nil || !token.Valid {
			return "", err
		}

		claims, ok := token.Claims.(*data.UserClaims)
		if !ok {
			return "", err
		}

		editor, err := db.GetUsernameFromId(claims.UserID)
		if err != nil {
			return "", err
		}
		return editor, nil
	}
	return strings.Split(r.RemoteAddr, ":")[0], nil
}

func SetupEditRoute(rt *chi.Mux, db data.Datastore, sc *appsignals.SignalConnector) {
	// Setup subrouter for wiki editing
	rt.Route("/e", func(editrouter chi.Router) {
		editrouter.Route("/{title}", func(pagerouter chi.Router) {
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

				// Get editor name
				editor, err := GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
				}

				// Handle request
				newPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&newPage); err != nil {
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with parsing json request: %s", err)
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
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with inserting page into database: %s", err)
					return
				}

				// Make response
				resp := make(map[string]string)
				resp["pageTitle"] = newPage.Title

				jsonresp.JsonOK(w, resp, "Page created!")

				// Fire event
				sc.Fire("onPageCreate", [2]string{editor, newPage.Title})
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

				// Get editor name
				editor, err := GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
				}

				// Handle request
				uPage := new(data.Page)
				if err := json.NewDecoder(r.Body).Decode(&uPage); err != nil {
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with parsing json request: %s", err)
					return
				}
				// Add title if nil
				if uPage.Title == "" {
					uPage.Title = titleParam
				}
				// Lowercase page title
				uPage.Title = strings.ToLower(uPage.Title)

				// Update page from database
				if err := db.UpdatePage(uPage, editor); err != nil {
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with inserting page into database: %s", err)
					return
				}

				// Make response
				jsonresp.JsonOK(w, make(map[string]string), "Page updated!")

				// Fire event
				sc.Fire("onPageUpdate", [1]string{editor})
			}))

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

				// Get editor name
				editor, err := GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
				}
				// Check if creating an account if possible
				userGroups := userutils.GetUserGroups(r.RemoteAddr)
				proceed := false
				for _, v := range userGroups {
					if wikiconfig.UserGroups[v]["delete"] {
						proceed = true
						break
					}
				}
				if !proceed {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with deleting page: Permission denied", nil)
					return
				}
				// Handle request
				pageTitle := new(struct {
					Title string `json:"title"`
				})
				if err := json.NewDecoder(r.Body).Decode(&pageTitle); err != nil {
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with parsing json request: %s", err)
					return
				}
				// Delete page from database
				if err := db.DeletePage(pageTitle.Title); err != nil {
					jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error deleting page from database: %s", err)
					return
				}

				// Make response
				jsonresp.JsonOK(w, make(map[string]string), "Page deleted!")

				// Fire event
				sc.Fire("onPageDelete", [1]string{editor})
			})
		})
	})
}
