package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func getUser(username string) (user User, err error) {
	err = db.One("Username", username, &user)

	if user.ID == 0 {
		return user, fmt.Errorf("not found: %s", username)
	}
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
