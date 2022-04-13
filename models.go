package main

type SEX string

const (
	Male   SEX = "Male"
	Female SEX = "Female"
)

type User struct {
	FirstName string
	LastName  string
	Age       int16
	Sex       SEX
	Hoobies   string
	City      string
}
