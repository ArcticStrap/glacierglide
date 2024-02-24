package edit

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ChaosIsFramecode/horinezumi/appsignals"
	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/ChaosIsFramecode/horinezumi/utils/userutils"
	"github.com/go-chi/chi/v5"
)

// Authenticate token, default to ip address if empty token.

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
				editor, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
          return
				}

        // Check if creating a page is possible
        userGroups := db.GetUserGroups(editor)
        proceed := userutils.UserCan("create",userGroups)
        
				if !proceed {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with creating page: Permission denied", nil)
					return
				}

        // Check if user is suspended
        proceed,err = db.IsSuspended(editor)
        if err != nil {
          jsonresp.JsonERR(w,http.StatusBadRequest,"Error with creating page: %s",err)
          return
        }
        if !proceed {
          jsonresp.JsonERR(w,http.StatusUnauthorized,"This user is temprorarily suspended.",nil)
          return
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
				editor, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
          return
				}

        // Check if editing a page is possible
        userGroups := db.GetUserGroups(editor)
        proceed := userutils.UserCan("edit",userGroups)

				if !proceed {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with editing page: Permission denied", nil)
					return
				}

        // Check if editor has sufficient permissions to edit
        minGroup, err := db.GetLockStatus(titleParam)
        if err != nil {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with editing page: %s", err)
          return
        }

        proceed = false
        for i := 0; i < len(userGroups); i++ {
          if i >= minGroup {
            proceed = true
            break
          }
        }
        if !proceed {
          jsonresp.JsonERR(w,http.StatusUnauthorized,"This page has been locked. User has insufficient permissions to edit",nil)
          return
        }
        
        // Check if user is suspended
        proceed,err = db.IsSuspended(editor)
        if err != nil {
          jsonresp.JsonERR(w,http.StatusBadRequest,"Error with editing page: %s",err)
          return
        }
        if !proceed {
          jsonresp.JsonERR(w,http.StatusUnauthorized,"This user is temprorarily suspended.",nil)
          return
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
				editor, err := data.GetLoginStatus(r.Header.Get("authtoken"), r, db)
				if err != nil {
					jsonresp.JsonERR(w, http.StatusBadRequest, "Error with authenticating user: %s", err)
          return
				}
				// Check if deleting a page is possible
				userGroups := db.GetUserGroups(editor)
				proceed := userutils.UserCan("delete",userGroups)
        if !proceed {
					jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with deleting page: Permission denied", nil)
					return
				}
        
        // Check if user is suspended
        proceed,err = db.IsSuspended(editor)
        if err != nil {
          jsonresp.JsonERR(w,http.StatusBadRequest,"Error with deleting page: %s",err)
          return
        }
        if !proceed {
          jsonresp.JsonERR(w,http.StatusUnauthorized,"This user is temprorarily suspended.",nil)
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
