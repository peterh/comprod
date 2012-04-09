package main

import (
	"html/template"
	"log"
	"net/http"
)

type handler struct {
	t *template.Template
}

func login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/login.html", 307)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(name) >= 1 {
		http.SetCookie(w, &http.Cookie{Name: "name", Value: name})
	} else {
		c, err := r.Cookie("name")
		if err != nil {
			login(w, r)
			return
		}
		name = c.Value
		if len(name) < 1 {
			login(w, r)
			return
		}
	}

	type data struct {
		Name string
	}
	h.t.Execute(w, &data{Name: name})
}

func main() {
	gameTemplate, err := template.ParseFiles("/Users/peterh/src/comprod/templates/game.html")
	if err != nil {
		log.Fatal("Fatal Error: ", err)
	}

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("/Users/peterh/src/comprod/static/"))))
	http.Handle("/", &handler{gameTemplate})

	log.Println("comprod started")
	log.Fatal(http.ListenAndServe(":2012", nil))
}

