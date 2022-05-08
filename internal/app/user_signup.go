package app

import "net/http"

func userSignup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	service.tmpl.ExecuteTemplate(w, "Signup", nil)
}
