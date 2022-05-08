package template

import (
	"html/template"
	"log"
)

func GetTemplate() *template.Template {
	log.Println("Templates OK")
	tmpl := template.Must(template.ParseGlob("form/*"))
	log.Println("Checking templates...")
	return tmpl
}
