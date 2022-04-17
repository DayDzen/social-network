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
)

type SEX string

const (
	Male   SEX = "Male"
	Female SEX = "Female"
	EMPTY  SEX = ""
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Age       int16
	Sex       SEX
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
	// dbConn := connectToDB()
	// defer dbConn.Close()

	selDB, err := dbConn.Query("SELECT * FROM user_profiles ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	res := []User{}
	for selDB.Next() {
		var id int64
		var age int16
		var fName, lName, sex, hobbies, city string
		err = selDB.Scan(&id, &fName, &lName, &age, &sex, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}
		user.ID = id
		user.FirstName = fName
		user.LastName = lName
		user.Age = age
		user.Sex = SEX(sex)
		user.Hobbies = hobbies
		user.City = city

		res = append(res, user)
	}
	log.Println(res)
	tmpl.ExecuteTemplate(w, "Index", res)
}

func Show(w http.ResponseWriter, r *http.Request) {
	// dbConn := connectToDB()
	// defer dbConn.Close()

	nId := r.URL.Query().Get("id")
	selDB, err := dbConn.Query("SELECT * FROM user_profiles WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	user := User{}
	for selDB.Next() {
		var id int64
		var age int16
		var fName, lName, sex, hobbies, city string
		err = selDB.Scan(&id, &fName, &lName, &age, &sex, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}

		user.ID = id
		user.FirstName = fName
		user.LastName = lName
		user.Age = age
		user.Sex = SEX(sex)
		user.Hobbies = hobbies
		user.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", user)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	// dbConn := connectToDB()
	// defer dbConn.Close()

	if r.Method == "POST" {
		fName := r.FormValue("first_name")
		lName := r.FormValue("last_name")
		age := r.FormValue("age")
		sex := r.FormValue("sex")
		hobbies := r.FormValue("hobbies")
		city := r.FormValue("city")
		insForm, err := dbConn.Prepare("INSERT INTO user_profiles(first_name, last_name, age, sex, hobbies, city) VALUES(?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fName, lName, age, sex, hobbies, city)
		log.Printf("INSERT:\nFirst Name: %v\nLast Name: %v\nAge: %v\nSex: %v\nHobbies: %v\nCity: %v\n", fName, lName, age, sex, hobbies, city)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/insert", Insert)

	http.ListenAndServe(":8080", nil)
}
