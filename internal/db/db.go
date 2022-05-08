package db

import (
	"fmt"
	"social-network/internal/model"
)

func GetAllUsers() ([]model.User, error) {
	query, err := dbConn.Query("SELECT id, first_name, last_name FROM users ORDER BY id DESC")
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers error: %w", err)
	}

	users := []model.User{}
	for query.Next() {
		var id int64
		var fName, lName string

		if err = query.Scan(&id, &fName, &lName); err != nil {
			return nil, fmt.Errorf("GetAllUsers error: %w", err)
		}

		users = append(users, model.User{
			ID:        id,
			FirstName: fName,
			LastName:  lName,
		})
	}

	return users, nil
}

func CreateNewUser(user model.User) error {
	insForm, err := dbConn.Prepare("INSERT INTO users(login, password, first_name, last_name, age, gender, hobbies, city) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		err = fmt.Errorf("SignUp err: %w", err)
		return err
	}

	if _, err = insForm.Exec(user.Login, user.Password, user.FirstName, user.LastName, user.Age, user.Gender, user.Hobbies, user.City); err != nil {
		err = fmt.Errorf("SignUp err: %w", err)
		return err
	}

	return nil
}

func GetUserPassByLogin(login string) (string, error) {
	selDB := dbConn.QueryRow("SELECT password FROM users WHERE login=?", login)
	if selDB.Err() != nil {
		return "", fmt.Errorf("GetPassByLogin err: %w", selDB.Err())
	}

	var dbPass string
	if err := selDB.Scan(&dbPass); err != nil {
		return "", fmt.Errorf("GetPassByLogin err: %w", err)
	}

	return dbPass, nil
}

func GetUserByID(id string) (model.User, error) {
	selDB, err := dbConn.Query("SELECT first_name, last_name, age, gender, hobbies, city FROM users WHERE id=?", id)
	if err != nil {
		return model.User{}, fmt.Errorf("GetUserByID err: %w", err)
		// panic(fmt.Errorf("UserProfile err: %w", err))
	}

	user := model.User{}
	for selDB.Next() {
		var age int
		var fName, lName, gender, hobbies, city string

		err = selDB.Scan(&fName, &lName, &age, &gender, &hobbies, &city)
		if err != nil {
			return model.User{}, fmt.Errorf("GetUserByID err: %w", err)
		}

		user.FirstName = fName
		user.LastName = lName
		user.Age = age
		user.Gender = model.GENDER(gender)
		user.Hobbies = hobbies
		user.City = city
	}

	return user, nil

}
