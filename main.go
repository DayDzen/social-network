package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type GENDER string

const (
	Male   GENDER = "Male"
	Female GENDER = "Female"
	EMPTY  GENDER = ""
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Age       int16
	Gender    GENDER
	Hobbies   string
	City      string
}

func connectToDB() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		return nil, fmt.Errorf("connectToDB err: %w", err)
	}

	return db, nil
}

var tmpl *template.Template
var dbConn *sql.DB

func checkEnvs() error {
	log.Println("Checking envs file...")
	var err error
	if err = godotenv.Load("prod_config.yaml"); err != nil {
		if err = godotenv.Load("local_config.yaml"); err != nil {
			return fmt.Errorf("checkEnvs err: %w", err)
		} else {
			log.Println("Local envs OK")
		}
	} else {
		log.Println("Prod envs OK")
	}

	return nil
}

func init() {
	var err error
	if err = checkEnvs(); err != nil {
		log.Print("No config file found")
	}

	log.Println("Checking DB connection...")
	dbConn, err = connectToDB()
	if err != nil {
		panic(err.Error())
	}
	log.Println("DB connection OK")

	log.Println("Checking templates...")
	tmpl = template.Must(template.ParseGlob("form/*"))
	log.Println("Templates OK")
}

func Index(w http.ResponseWriter, r *http.Request) {
	selDB, err := dbConn.Query("SELECT id, first_name, last_name FROM users ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	res := []User{}
	for selDB.Next() {
		var id int64
		var fName, lName string
		err = selDB.Scan(&id, &fName, &lName)
		// var age int16
		// var gender, hobbies, city string
		// err = selDB.Scan(&id, &fName, &lName, &age, &gender, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}
		user.ID = id
		user.FirstName = fName
		user.LastName = lName
		// user.Age = age
		// user.Gender = GENDER(gender)
		// user.Hobbies = hobbies
		// user.City = city

		res = append(res, user)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		login := r.FormValue("login")
		pass := getHash([]byte(r.FormValue("pass")))

		fName := r.FormValue("first_name")
		lName := r.FormValue("last_name")
		age := r.FormValue("age")
		gender := r.FormValue("gender")
		hobbies := r.FormValue("hobbies")
		city := r.FormValue("city")

		insForm, err := dbConn.Prepare("INSERT INTO users(login, password, first_name, last_name, age, gender, hobbies, city) VALUES(?,?,?,?,?,?,?,?)")
		if err != nil {
			err = fmt.Errorf("SignUp err: %w", err)
			panic(err.Error())
		}

		_, err = insForm.Exec(login, pass, fName, lName, age, gender, hobbies, city)
		if err != nil {
			err = fmt.Errorf("SignUp err: %w", err)
			panic(err.Error())
		}
		log.Printf("INSERT:\nFirst Name: %v\nLast Name: %v\nAge: %v\nGender: %v\nHobbies: %v\nCity: %v\n", fName, lName, age, gender, hobbies, city)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func Show(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := dbConn.Query("SELECT first_name, last_name, age, gender, hobbies, city FROM users WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	for selDB.Next() {
		var age int16
		var fName, lName, gender, hobbies, city string
		err = selDB.Scan(&fName, &lName, &age, &gender, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}
		user.FirstName = fName
		user.LastName = lName
		user.Age = age
		user.Gender = GENDER(gender)
		user.Hobbies = hobbies
		user.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", user)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

// func Insert(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		fName := r.FormValue("first_name")
// 		lName := r.FormValue("last_name")
// 		age := r.FormValue("age")
// 		gender := r.FormValue("gender")
// 		hobbies := r.FormValue("hobbies")
// 		city := r.FormValue("city")

// 		insForm, err := dbConn.Prepare("INSERT INTO users(first_name, last_name, age, gender, hobbies, city) VALUES(?,?,?,?,?,?)")
// 		if err != nil {
// 			panic(err.Error())
// 		}
// 		insForm.Exec(fName, lName, age, gender, hobbies, city)
// 		log.Printf("INSERT:\nFirst Name: %v\nLast Name: %v\nAge: %v\nGender: %v\nHobbies: %v\nCity: %v\n", fName, lName, age, gender, hobbies, city)
// 	}
// 	http.Redirect(w, r, "/", http.StatusMovedPermanently)
// }

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/", Index)

	http.HandleFunc("/sign-up", SignUp)
	// http.HandleFunc("/sign-in", SignIn)
	// http.HandleFunc("/sign-offss", SignOff)

	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	// http.HandleFunc("/insert", Insert)

	http.ListenAndServe(":8080", nil)
}
