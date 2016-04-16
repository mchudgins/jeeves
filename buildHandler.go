package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var buildRegex *regexp.Regexp

func init() {
	log.Printf("init called.")
	buildRegex = regexp.MustCompile("/builds/(?P<image>[a-zA-Z:]*)/(?P<buildName>[[:alpha:]]*)(?P<excess>/.*){0,}")
}

func BuildHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("processing http request: %v\n", *r)

	url := r.URL.RequestURI()
	log.Printf("URL:  %s\n", url)
	log.Printf("Query: %v", r.URL.Query())
	log.Printf("size: %d", len(buildRegex.FindStringSubmatch(url)))
	for _, s := range buildRegex.FindStringSubmatch(url) {
		log.Printf("element: %s", s)
	}
	elements := buildRegex.FindStringSubmatch(url)
	if len(elements) != 4 {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	imageName := elements[1]
	buildName := elements[2]

	log.Printf("imageName: %s; buildName: %s", imageName, buildName)

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

		http.Error(w, http.StatusText(400), 400)

	default:
		fmt.Fprintf(w, "%s : %q", r.Method, html.EscapeString(r.URL.Path))
	}
}
