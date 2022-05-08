package app

import "net/http"

func userLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	service.tmpl.ExecuteTemplate(w, "Login", nil)
}
