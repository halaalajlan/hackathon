package auth

import (
	"errors"
	"net/http"

	"github.com/halaalajlan/hackathon/models"
	"golang.org/x/crypto/bcrypt"
)

// ErrInvalidPassword is thrown when a user provides an incorrect password.
var ErrInvalidPassword = errors.New("Invalid Password")

// ErrPasswordMismatch is thrown when a user provides a blank password to the register
// or change password functions
var ErrPasswordMismatch = errors.New("Password cannot be blank")

// ErrEmptyPassword is thrown when a user provides a blank password to the register
// or change password functions
var ErrEmptyPassword = errors.New("No password provided")

// Login attempts to login the user given a request.
func Login(r *http.Request) (bool, models.Hospital, error) {
	username, password := r.FormValue("username"), r.FormValue("password")
	u, err := models.GetUserByUsername(username)
	if err != nil {
		return false, models.Hospital{}, err
	}
	//If we've made it here, we should have a valid user stored in u
	//Let's check the password
	err = bcrypt.CompareHashAndPassword([]byte(u.Hash), []byte(password))
	if err != nil {
		return false, models.Hospital{}, ErrInvalidPassword
	}
	return true, u, nil
}
