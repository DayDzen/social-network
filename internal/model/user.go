package model

type GENDER string

const (
	Male   GENDER = "Male"
	Female GENDER = "Female"
	EMPTY  GENDER = ""
)

type User struct {
	ID        int64
	Login     string
	Password  string
	FirstName string
	LastName  string
	Age       int
	Gender    GENDER
	Hobbies   string
	City      string
}

func NewUser(login, pass, firstName, lastName, hobbies, city string, age int, gender GENDER) User {
	return User{
		Login:     login,
		Password:  pass,
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
		Gender:    gender,
		Hobbies:   hobbies,
		City:      city,
	}
}
