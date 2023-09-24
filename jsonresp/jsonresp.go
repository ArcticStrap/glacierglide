package jsonresp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func JsonOK(w http.ResponseWriter, resp map[string]string, message string) {
	// Make response
	resp["message"] = message
	resp["code"] = strconv.Itoa(200)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func JsonERR(w http.ResponseWriter, code int, message string, msgerr error) {
	message = fmt.Sprintf(message, msgerr)
	// Make response
	resp := make(map[string]string)
	resp["message"] = message
	resp["code"] = strconv.Itoa(code)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}
