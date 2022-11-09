package rest

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, r *http.Request, code int, message string) {
	addStandardHeaders(w, r)
	respondWithJSON(w, r, code, map[string]string{"error": message})
}

func respondOK(w http.ResponseWriter, r *http.Request, code int) {
	addStandardHeaders(w, r)
	respondWithJSON(w, r, code, map[string]string{"result": "success"})
}

func respondEmpty(w http.ResponseWriter, r *http.Request, code int) {
	addStandardHeaders(w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.WriteHeader(code)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	addStandardHeaders(w, r)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	response, _ := json.Marshal(payload)
	w.Write(response)
}

func addStandardHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
