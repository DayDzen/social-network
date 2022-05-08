package app

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var service Serivce

type Serivce struct {
	tmpl *template.Template
}

func NewService(template *template.Template) {
	service.tmpl = template
}

func StartService(r *mux.Router) {
	r.HandleFunc("/", index)
	r.HandleFunc("/signup", userSignup)
	r.HandleFunc("/create-user", createUser).Methods(http.MethodPost)
	// http.HandleFunc("/sign-offss", SignOff)
	r.HandleFunc("/login", userLogin)
	r.HandleFunc("/auth-user", authUser).Methods(http.MethodPost)

	r.HandleFunc("/user-profile", userProfile).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", r))
}
