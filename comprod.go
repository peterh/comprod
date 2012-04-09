package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type handler struct {
	t *template.Template
	g *gameState
}

func login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/login.html", 307)
}

func formatValue(value uint64) template.HTML {
	s := strconv.FormatUint(value, 10)
	chunk := make([]string, 0)
	for len(s) > 0 {
		if len(s) >= 3 {
			chunk = append([]string{s[len(s)-3:]}, chunk...)
			s = s[:len(s)-3]
		} else {
			chunk = append([]string{s}, chunk...)
			s = ""
		}
	}
	return template.HTML(strings.Join(chunk, "&thinsp;"))
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

	type entry struct {
		Name   string
		Cost   uint64
		Shares uint64
		Value  template.HTML
	}
	type data struct {
		Name   string
		Stocks []entry
		Cash   template.HTML
	}
	s := h.g.listStocks()
	p := h.g.player(name)
	d := &data{Name: name}
	for k, v := range s {
		d.Stocks = append(d.Stocks, entry{
			Name:   v.Name,
			Cost:   v.Value,
			Shares: p.Shares[k],
			Value:  formatValue(p.Shares[k] * v.Value),
		})
	}
	d.Cash = formatValue(p.Cash)
	h.t.Execute(w, d)
}

func main() {
	gameTemplate, err := template.ParseFiles("/Users/peterh/src/comprod/templates/game.html")
	if err != nil {
		log.Fatal("Fatal Error: ", err)
	}

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("/Users/peterh/src/comprod/static/"))))
	http.Handle("/", &handler{gameTemplate, newGame()})

	log.Println("comprod started")
	log.Fatal(http.ListenAndServe(":2012", nil))
}

