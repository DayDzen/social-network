package app

import (
	"fmt"
	"net/http"
	"social-network/internal/db"
)

func index(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	users, err := db.GetAllUsers()
	if err != nil {
		panic(fmt.Errorf("index error: %w", err))
	}

	service.tmpl.ExecuteTemplate(w, "Index", users)
}
