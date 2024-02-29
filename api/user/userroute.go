package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ChaosIsFramecode/horinezumi/appsignals"
	"github.com/ChaosIsFramecode/horinezumi/data"
	"github.com/ChaosIsFramecode/horinezumi/jsonresp"
	"github.com/ChaosIsFramecode/horinezumi/utils/userutils"
	"github.com/go-chi/chi/v5"
)

func SetupUserRoute(rt chi.Router, db data.Datastore, sc *appsignals.SignalConnector) {
	rt.Post("/CreateAccount", func(w http.ResponseWriter, r *http.Request) {
		// Expect json response
		w.Header().Set("Content-Type", "application/json")

		// Check if creating an account is possible
		userGroups, err := db.GetUserGroups(r.RemoteAddr)
		if err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with getting user groups: %s", err)
			return
		}
		proceed := userutils.UserCan("createaccount", userGroups)

		if !proceed {
			jsonresp.JsonERR(w, http.StatusUnauthorized, "Error with creating user account: %s", fmt.Errorf("Permission denied"))
			return
		}

		// Decode account request
		var createReq data.AccountReq

		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with decoding json: ", err)
			return
		}

		// Prevent empty account creation
		if createReq.Username == "" || createReq.Password == "" {
			jsonresp.JsonERR(w, http.StatusBadRequest, "%s", fmt.Errorf("invalid credentials"))
			return
		}

		newUser, err := db.CreateUser(createReq.Username, createReq.Password)
		if err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with creating user account: ", err)
			return
		}

		resp := make(map[string]string)
		resp["password"] = createReq.Password
		resp["creation-date"] = newUser.CreationDate.Format("2006-01-02 15:04:05") + " UTC"

		jsonresp.JsonOK(w, resp, fmt.Sprintf("User account %s has been created", newUser.Username))

		// Fire event
		sc.Fire("onCreateAccount", [1]string{createReq.Username})
	})

	rt.Post("/Login", func(w http.ResponseWriter, r *http.Request) {
		// Expect json response
		w.Header().Set("Content-Type", "application/json")

		// Check if session exists
		_, err := r.Cookie("sessionauth")
		if err == nil {
			jsonresp.JsonERR(w, http.StatusUnauthorized, "User already has an active session", nil)
			return
		}

		// Decode account request
		var loginReq data.AccountReq

		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with decoding json: ", err)
			return
		}

		// Prevent empty account login
		if loginReq.Username == "" || loginReq.Password == "" {
			jsonresp.JsonERR(w, http.StatusBadRequest, "%s", fmt.Errorf("invalid credentials"))
			return
		}

		// Get user
		u, err := db.GetUser(loginReq.Username)
		if err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with fetching user info: ", err)
			return
		}

		// Check if password is valid
		if !u.ValidPassword(loginReq.Password) {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with logging into user: %s", fmt.Errorf("authentication error"))
			return
		}

		token, err := data.CreateJWT(u)
		if err != nil {
			jsonresp.JsonERR(w, http.StatusBadRequest, "Error with fetching auth token: ", err)
			return
		}

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionauth",
			Value:    token,
			HttpOnly: true,
		})

		resp := make(map[string]string)
		resp["token"] = token

		jsonresp.JsonOK(w, resp, fmt.Sprintf("Successfully logged in as user %s.", loginReq.Username))
	})
	rt.Post("/Logout", func(w http.ResponseWriter, _ *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "sessionauth",
			Value:    "",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		})

		resp := make(map[string]string)

		jsonresp.JsonOK(w, resp, "Logged out")
	})
}
