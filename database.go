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

type Request struct {
	Query   string
	Respond Table
	Done    chan error
}

type Task struct {
	Question string
	Solution string
	Respond  Table
}

type Database struct {
	Addr string
	User string
	Pass string
	Name string
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

func (d *Database) runQuerys(req chan *Request) error {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s", d.User, d.Name, d.Pass)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	queue := make(chan int, QueueRate)
	for r := range req {
		queue <- 1

		go func(r *Request) {
			r.Done <- r.query(db)
			<-queue
		}(r)
	}

	return nil
}
