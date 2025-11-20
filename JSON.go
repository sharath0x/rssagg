package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	if code > 499 {
		log.Printf("Respond with 5xx error %v", msg)
	}

	type Error struct {
		Error string `json:"error"`
	}
	err := Error{msg}
	respondWithJson(w, code, err)
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	// code - statuscode
	// payload - marshal to json

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error occured on server side: %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
