package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/halaalajlan/hackathon/logger"
	"github.com/halaalajlan/hackathon/models"

	"golang.org/x/crypto/bcrypt"
)

func (as *Server) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		JSONResponse(w, models.Response{Success: false, Message: "Method not allowed"}, http.StatusBadRequest)
		return
	}
	username, password := r.FormValue("username"), r.FormValue("password")
	u, err := models.GetUserByUsername(username)
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusNotFound)
		return
	}
	//If we've made it here, we should have a valid user stored in u
	//Let's check the password
	err = bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		JSONResponse(w, models.Response{Success: false, Message: "Invalid Password"}, http.StatusUnauthorized)
		return
	}
	JSONResponse(w, u.ApiKey, http.StatusOK)

}

func (as *Server) Patients(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

	case "POST":

	}
}

// JSONResponse attempts to set the status code, c, and marshal the given interface, d, into a response that
// is written to the given ResponseWriter.
func JSONResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		log.Error(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}
