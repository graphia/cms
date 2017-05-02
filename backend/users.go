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

func createUser(user User) (err error) {
	user.Password, err = bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost)

	err = db.Save(&user)
	if err != nil {
		return fmt.Errorf("User not created, %s", err.Error())
	}
	return nil
}

func allUsers() (users []User, err error) {
	err = db.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
