package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func getUserByID(id int) (user User, err error) {
	err = db.One("ID", id, &user)

	Debug.Println("Finding user by ID", id)

	if user.ID == 0 {
		Debug.Println("Cannot find user with ID", id)

		return user, fmt.Errorf("not found: %d", id)
	}
	Debug.Println("Found user", id)

	return user, err
}

func getUserByUsername(username string) (user User, err error) {
	err = db.One("Username", username, &user)
	Debug.Println("Finding user by username", username)

	if user.ID == 0 {
		Debug.Println("Cannot find user with Username", username)

		return user, fmt.Errorf("not found: %s", username)
	}

	Debug.Println("Found user", username)

	return user, err
}

func getLimitedUserByUsername(username string) (limitedUser LimitedUser, err error) {
	var user User
	err = db.One("Username", username, &user)
	Debug.Println("Finding user by username", username)

	if user.ID == 0 {
		Debug.Println("Cannot find user with Username", username)

		return limitedUser, fmt.Errorf("not found: %s", username)
	}

	Debug.Println("Found user", username)
	limitedUser = convertToLimitedUser(user)

	return limitedUser, err
}

func createUser(user User) (err error) {

	user.Password, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)

	err = validate.Struct(user)
	if err != nil {
		return err
	}

	Debug.Println("Validation passed, saving")

	err = db.Save(&user)
	if err != nil {
		return fmt.Errorf("User not created, %s", err.Error())
	}
	return nil
}

func allUsers() (limitedUsers []LimitedUser, err error) {

	var users []User

	err = db.All(&users)

	for _, user := range users {
		limitedUsers = append(limitedUsers, convertToLimitedUser(user))
	}

	if err != nil {
		return nil, err
	}
	return limitedUsers, nil
}

func countUsers() (qty int, err error) {
	qty, err = db.Count(&User{})
	if err != nil {
		return -1, err
	}
	return qty, err
}

func convertToLimitedUser(user User) LimitedUser {
	return LimitedUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Name,
	}
}

// Actually do a 'hard' delete
func deleteUser(user User) bool {
	err := db.DeleteStruct(&user)
	if err != nil {
		// handle
	}
	return true
}

// Soft delete, prevent user from signing in and
// accessing the API
func deactivateUser(user User) error {
	err := db.UpdateField(&user, "Active", false)
	if err != nil {
		return err
	}
	return nil
}

func reactivateUser(user User) error {
	err := db.UpdateField(&user, "Active", true)
	if err != nil {
		return err
	}
	return nil
}
