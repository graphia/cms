package main

import "fmt"

func getUser(username string) (user User, err error) {
	err = db.One("Username", username, &user)

	if user.ID == 0 {
		return user, fmt.Errorf("not found: %s", username)
	}
	return user, err
}

func createUser(user User) (User, error) {
	err := db.Save(&user)
	return user, err
}
