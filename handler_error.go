package main

import "net/http"

func handler_error(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, 400, "something went wrong")
}
