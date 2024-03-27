package search

import (
	"encoding/json"
	"net/http"

	"github.com/ArcticStrap/glacierglide/data"
	"github.com/ArcticStrap/glacierglide/jsonresp"
)

type SearchRequest struct {
	Pattern string `json:"pattern"`
	Limit   int    `json:"limit"`
}

func SetupSearchEngine(rt *http.ServeMux, db data.Datastore) {
	rt.HandleFunc("POST /api/search", func(w http.ResponseWriter, r *http.Request) {
		// Expect json response
		w.Header().Set("Content-Type", "application/json")

		// Decode body
		var sr SearchRequest
		if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
			jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with parsing json request: %s", err)
			return
		}

		if sr.Pattern == "" {
			return
		}

		// Check whether search limit is between 1 and 500.
		if sr.Limit < 1 || sr.Limit > 500 {
			sr.Limit = 10
		}

		res, err := db.SearchPagesFromTitlePrefix(sr.Pattern, sr.Limit)
		if err != nil {
			jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with searching pages: %s", err)
			return
		}

		if len(res) == 0 {
			res, err = db.SearchPagesContainingTitle(sr.Pattern, sr.Limit)
			if err != nil {
				jsonresp.JsonERR(w, http.StatusUnprocessableEntity, "Error with searching pages: %s", err)
				return
			}
		}

		var matches []string
		for _, v := range res {
			matches = append(matches, v.Title)
		}

		// Return possible matches
		json.NewEncoder(w).Encode(matches)
	})
}
