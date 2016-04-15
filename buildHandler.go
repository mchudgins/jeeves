package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func BuildHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("processing htp request: %v\n", *r)

	switch strings.ToUpper(r.Method) {
	case "POST":
		if r.Body != nil {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Fatal(err)

			}
			defer r.Body.Close()

			if body != nil && len(body) > 0 {
				log.Printf("read: %s\n\n", body)
			}
		}

	default:
		fmt.Fprintf(w, "%s : %q", r.Method, html.EscapeString(r.URL.Path))
	}
}
