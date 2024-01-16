package main

import (
	"log"
	"net/http"
)

var (
	querys chan *Request
	tasks  map[string]*Task
)

func main() {
	var err error
	db := &Database{
		Addr: "",
		User: "",
		Pass: "",
		Name: "",
	}

	querys = make(chan *Request)
	go func() {
		err = db.runQuerys(querys)
		if err != nil {
			log.Fatal(err)
		}
	}()

	tasks, err = loadTasks("tasks.json")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/task/", handleTask)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
