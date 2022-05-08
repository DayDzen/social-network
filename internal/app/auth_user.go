package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"social-network/internal/db"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func authUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	login := strings.TrimSpace(r.FormValue("login"))
	userPass := strings.TrimSpace(r.FormValue("pass"))

	dbPass, err := db.GetUserPassByLogin(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(fmt.Errorf("authUser err: %w", err))
			w.Write([]byte(`{"response":"No user with this login!"}`))
			return
		}

		log.Panic(fmt.Errorf("authUser err: %w", err))
	}

	userPassByte := []byte(userPass)
	dbPassByte := []byte(dbPass)

	if passErr := bcrypt.CompareHashAndPassword(dbPassByte, userPassByte); passErr != nil {
		w.Write([]byte(`{"response":"Wrong Password!"}`))
		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
