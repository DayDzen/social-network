package main

import (
	"database/sql"
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

type UserProfile struct {
	ID        int64
	FirstName string
	LastName  string
	Age       int16
	Sex       SEX
	Hobbies   string
	City      string
}

func connectToDB() *sql.DB {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}

var tmpl *template.Template

func init() {
	log.Println("INIT Started")

	if err := godotenv.Load("prod_config.yaml"); err != nil {
		if err = godotenv.Load("local_config.yaml"); err != nil {
			log.Print("No config file found")
		}
	}

	tmpl = template.Must(template.ParseGlob("form/*"))
	log.Println("SUCCESS init template")
}

func Index(w http.ResponseWriter, r *http.Request) {
	dbConn := connectToDB()
	defer dbConn.Close()

	selDB, err := dbConn.Query("SELECT * FROM user_profiles ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}
	uProfile := UserProfile{}
	res := []UserProfile{}
	for selDB.Next() {
		var id int64
		var age int16
		var fName, lName, sex, hobbies, city string
		err = selDB.Scan(&id, &fName, &lName, &age, &sex, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}
		uProfile.ID = id
		uProfile.FirstName = fName
		uProfile.LastName = lName
		uProfile.Age = age
		uProfile.Sex = SEX(sex)
		uProfile.Hobbies = hobbies
		uProfile.City = city

		res = append(res, uProfile)
	}
	tmpl.ExecuteTemplate(w, "Index", res)
}

func Show(w http.ResponseWriter, r *http.Request) {
	dbConn := connectToDB()
	defer dbConn.Close()

	nId := r.URL.Query().Get("id")
	selDB, err := dbConn.Query("SELECT * FROM user_profiles WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}
	uProfile := UserProfile{}
	for selDB.Next() {
		var id int64
		var age int16
		var fName, lName, sex, hobbies, city string
		err = selDB.Scan(&id, &fName, &lName, &age, &sex, &hobbies, &city)
		if err != nil {
			panic(err.Error())
		}

		uProfile.ID = id
		uProfile.FirstName = fName
		uProfile.LastName = lName
		uProfile.Age = age
		uProfile.Sex = SEX(sex)
		uProfile.Hobbies = hobbies
		uProfile.City = city
	}
	tmpl.ExecuteTemplate(w, "Show", uProfile)
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Insert(w http.ResponseWriter, r *http.Request) {
	dbConn := connectToDB()
	defer dbConn.Close()

	if r.Method == "POST" {
		fName := r.FormValue("f_name")
		lName := r.FormValue("l_name")
		age := r.FormValue("age")
		sex := r.FormValue("sex")
		hobbies := r.FormValue("hobbies")
		city := r.FormValue("city")
		insForm, err := dbConn.Prepare("INSERT INTO user_profiles(f_name, l_name, age, sex, hobbies, city) VALUES(?,?,?,?,?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(fName, lName, city)
		log.Printf("INSERT:\nfName: %v\nlName: %v\nAge: %v\nSex: %v\nHobbies: %v\nCity: %v\n", fName, lName, age, sex, hobbies, city)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/edit", Insert)

	http.ListenAndServe(":8080", nil)
}
