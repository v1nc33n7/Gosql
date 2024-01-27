package main

import (
	"log"
)

const (
	USER   = ""
	PASS   = ""
	DBNAME = ""
)

func main() {
	db, err := createDatabase(USER, PASS, DBNAME)
	if err != nil {
		log.Panic(err)
	}

	wait := make(chan bool)
	go func() {
		err = db.runQueue(wait)
		if err != nil {
			log.Panic(err)
		}
	}()
	<-wait

	err = db.loadTable("ny_house", "ny_house_tasks")
	if err != nil {
		log.Panic(err)
	}

	runHandlers(db)
}
