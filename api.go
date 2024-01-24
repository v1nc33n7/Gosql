package main

import (
	"fmt"
	"net/http"
	"regexp"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	re := regexp.MustCompile(`^/(\w+)/(\d)$`)
	fmt.Printf("%q\n", re.FindStringSubmatch(r.URL.Path))
}

func runHandlers(db *Database) {
	http.HandleFunc("/", handleIndex)

	http.ListenAndServe(":3000", nil)
}
