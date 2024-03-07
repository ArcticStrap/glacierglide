package jsonresp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Gives the http 200 OK response to a ResponseWriter
func JsonOK(w http.ResponseWriter, resp map[string]string, message string) {
	// Write OK CODE
	w.WriteHeader(http.StatusOK)

	// Make response
	resp["message"] = message

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Write(jsonResp)
}

// Gives an http error code to a ResponseWriter
func JsonERR(w http.ResponseWriter, code int, message string, msgerr error) {
	// Write ERR CODE
	w.WriteHeader(code)

	// Format error message
	message = fmt.Sprintf(message, msgerr)
	// Make response
	resp := make(map[string]string)
	resp["message"] = message

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
		return
	}

	w.Write(jsonResp)
}

// Error handler for cookies
func FetchCookieValue(name string, r *http.Request) (string,error) {
  cookie, err := r.Cookie(name)
  if err != nil {
    if errors.Is(err,http.ErrNoCookie) {
      return "", nil
    }
    return "", err
  }

  return cookie.Value, nil
}
