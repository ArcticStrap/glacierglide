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

		createReq := new(struct {
			Username string `json:"username"`
			Password string `json:"password"`
		})

		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			jsonresp.JsonERR(w, 400, "Error with decoding json: ", err)
			return
		}

		newUser, err := db.CreateUser(createReq.Username, createReq.Password)
		if err != nil {
			jsonresp.JsonERR(w, 400, "Error with creating user account: ", err)
		}

		resp := make(map[string]string)
		resp["password"] = newUser.Password
		resp["creation-date"] = newUser.CreationDate.Format("2006-01-02 15:04:05") + " UTC"

		jsonresp.JsonOK(w, resp, fmt.Sprintf("User account %s has been created", newUser.Username))
	})
}
