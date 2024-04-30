package server

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseGlob("web/landing/*.html"))

type Landing struct{}

func landing() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := Landing{}
		tmpl.ExecuteTemplate(w, "index.html", data)
	})
}
