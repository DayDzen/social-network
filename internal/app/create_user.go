package app

import (
	"fmt"
	"log"
	"net/http"
	"social-network/internal/db"
	"social-network/internal/model"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	login := strings.TrimSpace(r.FormValue("login"))
	pass := getHash([]byte(strings.TrimSpace(r.FormValue("pass"))))

	fName := strings.TrimSpace(r.FormValue("first_name"))
	lName := strings.TrimSpace(r.FormValue("last_name"))
	age, err := strconv.Atoi(strings.TrimSpace(r.FormValue("age")))
	if err != nil {
		err = fmt.Errorf("createUser err: %w", err)
		log.Println(err.Error())
	}
	gender := strings.TrimSpace(r.FormValue("gender"))
	hobbies := strings.TrimSpace(r.FormValue("hobbies"))
	city := strings.TrimSpace(r.FormValue("city"))

	user := model.NewUser(login, pass, fName, lName, hobbies, city, age, model.GENDER(gender))

	if err := db.CreateNewUser(user); err != nil {
		log.Panic(fmt.Errorf("createUser err: %w", err))
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
