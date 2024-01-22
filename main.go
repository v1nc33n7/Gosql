package main

import (
	"log"
	"net/http"
)

var (
	querys chan *Request
	tasks  map[string]*Task
)

func loadTasks(table string) error {
	req := &Request{
		Query: "SELECT * FROM " + table,
		Done:  make(chan error),
	}
	querys <- req

	err := <-req.Done
	if err != nil {
		return err
	}

	tasks = make(map[string]*Task)
	for k := range req.Respond.Rows {
		solution := &Request{
			Query: req.Respond.Rows[k][2],
			Done:  make(chan error),
		}
		querys <- solution

		err = <-solution.Done
		if err != nil {
			return err
		}

		tasks[req.Respond.Rows[k][0]] = &Task{
			Question: req.Respond.Rows[k][1],
			Solution: req.Respond.Rows[k][2],
			Respond:  solution.Respond,
		}
	}

	return err
}

func main() {
	db := &Database{
		Addr: "",
		User: "",
		Pass: "",
		Name: "",
	}

	querys = make(chan *Request)
	go func() {
		err := db.runQuerys(querys)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := loadTasks("ny_house_tasks")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/task/", handleTask)
	http.HandleFunc("/anwser/", handleQuery)
	http.HandleFunc("/list", handleList)

	log.Printf("Server running on :3000\n")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
