package main

import (
	"log"
	"social-network/internal/app"
	"social-network/internal/config"
	"social-network/internal/db"
	"social-network/internal/template"

	"github.com/gorilla/mux"
)

func init() {
	if err := config.CheckEnvs(); err != nil {
		log.Panic("No config file found")
	}

	log.Println("Checking DB connection...")
	err := db.ConnectToDB()
	if err != nil {
		log.Panic(err.Error())
	}
	log.Println("DB connection OK")

	tmpl := template.GetTemplate()

	app.NewService(tmpl)
}

func main() {
	log.Println("Server started on: http://localhost:8080")

	router := mux.NewRouter()
	app.StartService(router)
}
