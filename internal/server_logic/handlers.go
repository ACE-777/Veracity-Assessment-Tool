package server_logic

import (
	"html/template"
	"log"
	"net/http"
)

const (
	firstRadio            = "First"
	secondRadio           = "Second"
	firstAlgorithmPython  = "scripts.algorithm_1_server"
	secondAlgorithmPython = "scripts.algorithm_2_server"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method must be GET", http.StatusNotFound)
	}

	if err := tmpl.ExecuteTemplate(w, "home.html", nil); err != nil {
		log.Fatalf("can not execute template login: %v", err)
	}

	w.WriteHeader(http.StatusOK)

	return
}

func coloring(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "method must be POST", http.StatusNotFound)
	}

	userInput := r.FormValue("textareainputOne")
	radio := r.FormValue("oneradio")
	if radio == firstRadio {
		getColoredFirst(userInput)
	}

	if radio == secondRadio {
		getColoredSecond(userInput)
	}

	var tmplTwo *template.Template

	t, err := tmplTwo.ParseFiles("internal/templates/coloring.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	if err := t.ExecuteTemplate(w, "coloring.html", nil); err != nil {
		log.Fatalf("can not execute template proceed: %v", err)
	}

	w.WriteHeader(http.StatusOK)

	return
}

func storeSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method must be POST", http.StatusNotFound)
	}

	userInput := r.FormValue("textareainputTwo")

	buildSearchDatabase(userInput)

	if err := tmpl.ExecuteTemplate(w, "store_sources.html", nil); err != nil {
		log.Fatalf("can not execute template proceed: %v", err)
	}

	w.WriteHeader(http.StatusOK)

	return
}
