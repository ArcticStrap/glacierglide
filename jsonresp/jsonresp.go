package jsonresp

import (
	"encoding/json"
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
