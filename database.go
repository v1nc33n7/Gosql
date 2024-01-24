package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

const QueueRate = 5

type Table struct {
	Columns []string
	Rows    [][]string
}

type Task struct {
	Question string
	Solution Table
}

type Request struct {
	Query   string
	Respond Table
	Done    chan error
}

type Database struct {
	// 1. List of all tables, tasks and solutions
	// 2. Channel for processing querys
	// 3. Postgres driver

	Tables map[string]map[string]*Task
	Listen chan *Request
	Driver *sql.DB
}

func (r *Request) query(db *sql.DB) error {
	rows, err := db.Query(r.Query)
	if err != nil {
		return err
	}

	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for _, c := range cols {
		r.Respond.Columns = append(r.Respond.Columns, strings.ToUpper(c))
	}

	row := make([][]byte, len(cols))
	anys := make([]any, len(cols))

	for i := range row {
		anys[i] = &row[i]
	}

	qResult := make([][]string, 0)
	for rows.Next() {
		err := rows.Scan(anys...)
		if err != nil {
			return err
		}

		sRow := make([]string, 0)
		for i := range row {
			sRow = append(sRow, string(row[i]))
		}

		qResult = append(qResult, sRow)
	}
	r.Respond.Rows = qResult

	return nil
}

func (d *Database) runQueue(wait chan bool) error {
	defer d.Driver.Close()
	defer close(d.Listen)

	queue := make(chan int, QueueRate)
	wait <- true
	for r := range d.Listen {
		queue <- 1

		go func(r *Request) {
			r.Done <- r.query(d.Driver)
			<-queue
		}(r)
	}

	return nil
}

func (d *Database) loadTable(table string, tasks string) error {
	reqTasks := &Request{
		Query: fmt.Sprintf("SELECT * FROM %s", tasks),
		Done:  make(chan error),
	}
	d.Listen <- reqTasks

	err := <-reqTasks.Done
	if err != nil {
		return err
	}

	d.Tables[table] = make(map[string]*Task)
	for k := range reqTasks.Respond.Rows {
		index := reqTasks.Respond.Rows[k][0]
		question := reqTasks.Respond.Rows[k][1]
		solution := reqTasks.Respond.Rows[k][2]

		reqSolution := &Request{
			Query: solution,
			Done:  make(chan error),
		}
		d.Listen <- reqSolution

		err := <-reqSolution.Done
		if err != nil {
			return err
		}

		d.Tables[table][index] = &Task{
			Question: question,
			Solution: reqSolution.Respond,
		}
	}

	return nil
}

func createDatabase(u string, p string, n string) (*Database, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s", u, p, n)

	pq, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db := &Database{
		Tables: make(map[string]map[string]*Task),
		Driver: pq,
		Listen: make(chan *Request),
	}

	return db, nil
}
