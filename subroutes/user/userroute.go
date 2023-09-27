package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/go-chi/chi/v5"
)

func SetupUserRoute(rt *chi.Mux, db data.Datastore) {
	rt.Post("/CreateAccount", func(w http.ResponseWriter, r *http.Request) {
		// Expect json response
		w.Header().Set("Content-Type", "application/json")

		// Decode account request
		var createReq data.AccountReq

		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			jsonresp.JsonERR(w, 400, "Error with decoding json: ", err)
			return
		}

		newUser, err := db.CreateUser(createReq.Username, createReq.Password)
		if err != nil {
			jsonresp.JsonERR(w, 400, "Error with creating user account: ", err)
		}

		resp := make(map[string]string)
		resp["password"] = createReq.Password
		resp["creation-date"] = newUser.CreationDate.Format("2006-01-02 15:04:05") + " UTC"

		jsonresp.JsonOK(w, resp, fmt.Sprintf("User account %s has been created", newUser.Username))
	})

	rt.Post("/Login", func(w http.ResponseWriter, r *http.Request) {
		// Expect json response
		w.Header().Set("Content-Type", "application/json")

		// Decode account request
		var loginReq data.AccountReq

		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			jsonresp.JsonERR(w, 400, "Error with decoding json: ", err)
			return
		}

		// Get user
		u, err := db.GetUser(loginReq.Username)
		if err != nil {
			jsonresp.JsonERR(w, 400, "Error with fetching user info: ", err)
			return
		}

		// Check if password is valid
		if !u.ValidPassword(loginReq.Password) {
			jsonresp.JsonERR(w, 400, "Error with logging into user: %s", fmt.Errorf("authentication error"))
			return
		}

		token, err := data.CreateJWT(u)
		if err != nil {
			jsonresp.JsonERR(w, 400, "Error with fetching auth token: ", err)
			return
		}

		resp := make(map[string]string)
		resp["token"] = token

		jsonresp.JsonOK(w, resp, fmt.Sprintf("Successfully logged in as user %s.", loginReq.Username))
	})
}
