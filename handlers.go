package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type Page struct {
	Number   int
	Question string

	Next     int
	Previous int
}

func urlNumber(r *http.Request) (int, string) {
	sNum := r.URL.Path[len("/task/"):]
	num, _ := strconv.Atoi(sNum)

	return num, sNum
}

func renderTemplate(w http.ResponseWriter, fn string, p *Page) {
	tmpl, _ := template.ParseFiles("tmpl/" + fn + ".html")
	tmpl.Execute(w, p)
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	num, sNum := urlNumber(r)
	if num == 0 || num > len(tasks) || num < 1 {
		http.Redirect(w, r, "/task/1", http.StatusFound)
	}

	task, ok := tasks[sNum]
	if !ok {
		return
	}

	p := &Page{
		Number:   num,
		Question: task.Solution,
		Next:     num + 1,
		Previous: num - 1,
	}
	renderTemplate(w, "task", p)
}
