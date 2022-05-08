package app

import (
	"fmt"
	"log"
	"net/http"
	"social-network/internal/db"
)

func userProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id := r.URL.Query().Get("id")

	user, err := db.GetUserByID(id)
	if err != nil {
		log.Panic(fmt.Errorf("userProfile err: %w", err))
	}

	service.tmpl.ExecuteTemplate(w, "UserProfile", user)
}
