package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var SECRET_KEY = []byte("gosecretkey")

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
		err = fmt.Errorf("Index err: %w", err)
		panic(err.Error())
	}

	res := []User{}
	for selDB.Next() {
		var id int64
		var fName, lName string

		if err = selDB.Scan(&id, &fName, &lName); err != nil {
			err = fmt.Errorf("Index err: %w", err)
			panic(err.Error())
		}

		res = append(res, User{
			ID:        id,
			FirstName: fName,
			LastName:  lName,
		})
	}
	tmpl.ExecuteTemplate(w, "Index", res)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
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

		if _, err = insForm.Exec(login, pass, fName, lName, age, gender, hobbies, city); err != nil {
			err = fmt.Errorf("SignUp err: %w", err)
			panic(err.Error())
		}

		log.Printf("INSERT:\nFirst Name: %v\nLast Name: %v\nAge: %v\nGender: %v\nHobbies: %v\nCity: %v\n", fName, lName, age, gender, hobbies, city)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func AuthUser(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	userPass := r.FormValue("pass")

	selDB := dbConn.QueryRow("SELECT password FROM users WHERE login=?", login)
	if selDB.Err() != nil {
		panic(fmt.Errorf("Show err: %w", selDB.Err()))
	}

	var dbPass string
	if err := selDB.Scan(&dbPass); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.Write([]byte(`{"response":"No user with this login!"}`))
			return
		}
		panic(fmt.Errorf("SignIn err: %w", err))
	}

	userPassByte := []byte(userPass)
	dbPassByte := []byte(dbPass)

	if passErr := bcrypt.CompareHashAndPassword(dbPassByte, userPassByte); passErr != nil {
		w.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	jwtToken, err := generateJWT()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	w.Write([]byte(`{"token":"` + jwtToken + `"}`))
}

func Show(w http.ResponseWriter, r *http.Request) {
	nId := r.URL.Query().Get("id")
	selDB, err := dbConn.Query("SELECT first_name, last_name, age, gender, hobbies, city FROM users WHERE id=?", nId)
	if err != nil {
		err = fmt.Errorf("Show err: %w", err)
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

func SignUp(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "SignUp", nil)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "SignIn", nil)
}

func getHash(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(SECRET_KEY)
	if err != nil {
		log.Println("Error in JWT token generation")
		return "", err
	}
	return tokenString, nil
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	http.HandleFunc("/", Index)

	http.HandleFunc("/sign-up", SignUp)
	http.HandleFunc("/sign-in", SignIn)
	// http.HandleFunc("/sign-offss", SignOff)

	http.HandleFunc("/create-user", CreateUser)
	http.HandleFunc("/auth-user", AuthUser)
	http.HandleFunc("/show", Show)

	http.ListenAndServe(":8080", nil)
}
