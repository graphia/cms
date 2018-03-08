package main

import (
	"encoding/json"
	"net/http"
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

	sr = SuccessResponse{
		Message: "User created",
	}

	JSONResponse(sr, http.StatusCreated, w)

}

func apiUpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	Debug.Println("req received")
	sr := SuccessResponse{Message: "User updated"}
	JSONResponse(sr, http.StatusCreated, w)
}

func apiDeleteUserHandler(w http.ResponseWriter, r *http.Request) {}
