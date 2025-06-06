package terminal

import (
	"html/template"
	"net/http"
)

func Show(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("terminal/terminal.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
