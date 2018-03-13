package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asdine/storm"

	"github.com/husobee/vestigo"
	validator "gopkg.in/go-playground/validator.v9"
)

// user admin functionality 👩🏽‍💻

// apiCreateUser
func apiCreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var sr SuccessResponse
	var fr FailureResponse
	var verrs map[string]string

	json.NewDecoder(r.Body).Decode(&user)

	user.Active = true

	err := createUser(user)
	verrs, err = checkUserModificationErrors(err)

	if len(verrs) > 0 {
		JSONResponse(verrs, http.StatusBadRequest, w)
		return
	}

	if err != nil {
		Error.Println("Could not update user", err.Error())
		fr = FailureResponse{Message: "User could not be updated"}
		JSONResponse(fr, http.StatusBadRequest, w)
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
	var verrs map[string]string

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
		return
	}

	err = updateUser(user, updates)
	verrs, err = checkUserModificationErrors(err)

	Debug.Println("verrs", verrs)

	if len(verrs) > 0 {
		Debug.Println("returning validation errors")
		JSONResponse(verrs, http.StatusBadRequest, w)
		return
	}

	if err != nil {
		Error.Println("Could not update user", err.Error())
		fr = FailureResponse{Message: "User could not be updated"}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	// success
	sr = SuccessResponse{Message: "User updated"}
	JSONResponse(sr, http.StatusOK, w)
}

func checkUserModificationErrors(errIn error) (verrs map[string]string, errOut error) {
	verrs = make(map[string]string)

	// skip unless we have *some* error
	if errIn != nil {

		switch errIn.(type) {
		// check if the error is a validation error
		case validator.ValidationErrors:
			Debug.Println("validation error", errIn)
			verrs = validationErrorsToJSON(errIn)
			return verrs, nil

		default:

			// another check, this time for storm's uniqueness test
			if errIn == storm.ErrAlreadyExists {
				Debug.Println("duplicate found, abort")
				verrs = make(map[string]string)
				verrs["Base"] = "is not unique"
				return verrs, nil
			}

			// finally, if it's not a validation or unique error, return it
			return nil, errIn
		}

	}
	return verrs, nil
}

func apiDeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var fr FailureResponse
	var sr SuccessResponse
	var user, currentUser User
	var err error

	username := vestigo.Param(r, "username")

	user, err = getUserByUsername(username)
	if err != nil {
		Error.Println("Failed to find user", username)
		fr = FailureResponse{Message: fmt.Sprintf("No user %s", username)}
		JSONResponse(fr, http.StatusNotFound, w)
		return
	}

	currentUser = getCurrentUser(r.Context())
	if user == currentUser {
		Warning.Printf("%s tried to delete themself", currentUser.Username)
		fr = FailureResponse{Message: "You cannot delete yourself"}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	err = user.delete()
	if err != nil {
		Error.Println("Couldn't delete user", username, err.Error())
		fr = FailureResponse{Message: fmt.Sprintf("Cannot delete user %s", username)}
		JSONResponse(fr, http.StatusBadRequest, w)
		return
	}

	sr = SuccessResponse{Message: fmt.Sprintf("Successfully deleted user %s", username)}
	JSONResponse(sr, http.StatusOK, w)

}

func apiActivateUserHandler(w http.ResponseWriter, r *http.Request)   {}
func apiDeactivateUserHandler(w http.ResponseWriter, r *http.Request) {}
