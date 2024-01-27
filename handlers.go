package main

import (
	"errors"
	"html/template"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
)

type Page struct {
	Name     string
	Number   string
	Question string
}

type Anwser struct {
	Respond Table
	Verify  bool
	Error   string
}

type List struct {
	Name  string
	Tasks []int
}

func findUrlVars(r *http.Request) (string, string, error) {
	re := regexp.MustCompile(`^/[a|q]/(\w+)/(\d)$`)
	urlVars := re.FindStringSubmatch(r.URL.Path)

	if len(urlVars) != 3 {
		return "", "", errors.New("wrong url")
	}

	return urlVars[1], urlVars[2], nil
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	t, _ := template.ParseFiles("tmpl/basic.html", "tmpl/"+tmpl+".html")
	t.ExecuteTemplate(w, "basic", p)
}

func handleList(w http.ResponseWriter, r *http.Request, db *Database) {
	l := make([]List, 0)

	for k, v := range db.Tasks {
		list := List{Name: k}
		for kv := range v {
			kvNumber, _ := strconv.Atoi(kv)
			list.Tasks = append(list.Tasks, kvNumber)
		}
		sort.Ints(list.Tasks)
		l = append(l, list)
	}

	t, _ := template.ParseFiles("tmpl/list.html")
	t.Execute(w, l)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func handleAnwser(w http.ResponseWriter, r *http.Request, db *Database) {
	table, number, err := findUrlVars(r)
	if err != nil {
		return
	}

	task, ok := db.Tasks[table][number]
	if !ok {
		return
	}

	req := &Request{
		Query: r.FormValue("anwser"),
		Done:  make(chan error),
	}
	db.Listen <- req

	t, _ := template.ParseFiles("tmpl/table.html")

	err = <-req.Done
	if err != nil {
		a := &Anwser{
			Error: "Error: " + err.Error()[4:],
		}
		t.Execute(w, a)
		return
	}

	a := &Anwser{
		Respond: req.Respond,
		Verify:  reflect.DeepEqual(task.Solution, req.Respond),
	}

	t.Execute(w, a)
}

func handleQuestion(w http.ResponseWriter, r *http.Request, db *Database) {
	table, number, err := findUrlVars(r)
	if err != nil {
		return
	}

	task, ok := db.Tasks[table][number]
	if !ok {
		return
	}

	p := &Page{
		Name:     table,
		Number:   number,
		Question: task.Question,
	}

	renderTemplate(w, "task", p)
}

func runHandlers(db *Database) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		handleList(w, r, db)
	})
	http.HandleFunc("/q/", func(w http.ResponseWriter, r *http.Request) {
		handleQuestion(w, r, db)
	})
	http.HandleFunc("/a/", func(w http.ResponseWriter, r *http.Request) {
		handleAnwser(w, r, db)
	})

	http.ListenAndServe(":3000", nil)
}
