package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/husobee/vestigo"
)

// user admin functionality 👩🏽‍💻

// apiCreateUser
func apiCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var sr SuccessResponse

	json.NewDecoder(r.Body).Decode(&user)

	user.Active = true

	err := createUser(user)
	if err != nil {
		errors := validationErrorsToJSON(err)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	Info.Println("User was created successfully", user.Username)

	sr = SuccessResponse{Message: "User created"}
	JSONResponse(sr, http.StatusCreated, w)

}

func apiUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var sr SuccessResponse
	var fr FailureResponse
	var err error
	var user User
	var updates LimitedUser

	username := vestigo.Param(r, "username")

	user, err = getUserByUsername(username)
	if err != nil {
		Error.Println("Failed to find user", username)
		fr = FailureResponse{Message: fmt.Sprintf("No user %s", username)}
		JSONResponse(fr, http.StatusNotFound, w)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		Error.Println("Could not decode payload", r.Body)
		fr = FailureResponse{Message: "Failed to decode payload"}
		JSONResponse(fr, http.StatusBadRequest, w)
	}

	err = updateUser(user, updates)
	if err != nil {
		errors := validationErrorsToJSON(err)
		Debug.Println("errors", errors)
		JSONResponse(errors, http.StatusBadRequest, w)
		return
	}

	// success
	sr = SuccessResponse{Message: "User updated"}
	JSONResponse(sr, http.StatusOK, w)
}

func apiDeleteUserHandler(w http.ResponseWriter, r *http.Request) {}
