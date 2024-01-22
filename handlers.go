package main

import (
	"fmt"
	"html/template"
	"net/http"
	"reflect"
	"strconv"
)

type Page struct {
	Number   int
	Question string

	Respond Table
	Verify  bool

	List []int

	Next     int
	Previous int
}

func urlNumber(r *http.Request, url string) (int, string) {
	sNum := r.URL.Path[len(url):]
	num, _ := strconv.Atoi(sNum)

	return num, sNum
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	num, sNum := urlNumber(r, "/task/")
	if num == 0 || num > len(tasks) || num < 1 {
		http.Redirect(w, r, "/task/1", http.StatusFound)
	}

	task, ok := tasks[sNum]
	if !ok {
		return
	}

	p := &Page{
		Number:   num,
		Question: task.Question,
		Next:     num + 1,
		Previous: num - 1,
	}

	tmpl, _ := template.ParseFiles("tmpl/task.html", "tmpl/basic.html")
	tmpl.ExecuteTemplate(w, "basic", p)
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	_, sNum := urlNumber(r, "/anwser/")
	anwser := r.FormValue("anwser")

	task, ok := tasks[sNum]
	if !ok {
		return
	}

	req := &Request{
		Query: anwser,
		Done:  make(chan error),
	}
	querys <- req

	err := <-req.Done
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	p := &Page{
		Respond: req.Respond,
		Verify:  reflect.DeepEqual(req.Respond, task.Respond),
	}

	tmpl, _ := template.ParseFiles("tmpl/table.html")
	tmpl.Execute(w, p)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	list := make([]int, 0)
	for i := 1; i <= len(tasks); i++ {
		list = append(list, i)
	}

	p := &Page{
		List: list,
	}

	tmpl, _ := template.ParseFiles("tmpl/list.html")
	tmpl.Execute(w, p)
}
