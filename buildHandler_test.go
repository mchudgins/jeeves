package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mchudgins/jeeves/pkg/k8sClient"
)

func TestGet(t *testing.T) {
	cases := []struct {
		url    string
		body   string
		expect int
	}{
		{url: "/builds/golang:latest", body: "", expect: 404},
		{url: "/builds/golang:latest/fubar", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar/1", body: "{ 'gitURL' : 'this'}", expect: 200},
		{url: "/builds/golang:latest/fubar/gorf", body: "", expect: 404},
		{url: "/builds/golang:latest/fubar?gorf", body: "", expect: 400},
	}

	var failed bool

	k8sClient := client.NewClientOrDie()
	dao, err := NewDaoBuilds()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	handler := NewBuildHandler(k8sClient, dao)

	for _, c := range cases {
		req, err := http.NewRequest("GET", "http://localhost"+c.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")

		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != c.expect {
			t.Logf("Expected HTTP Status Code of %d.  Got %d.  URL = %s\n", c.expect, w.Code, c.url)
			failed = true
		}

		if failed {
			t.Fail()
		}
	}
}

func TestPost(t *testing.T) {
	cases := []struct {
		url    string
		body   string
		expect int
	}{
		{url: "/builds/golang:latest", body: "", expect: 404},
		{url: "/builds/golang:latest/fubar", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar", body: "{ 'gitURL' : 'this'}", expect: 200},
		{url: "/builds/golang:latest/fubar/gorf", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar?gorf", body: "", expect: 400},
	}

	var failed bool

	k8sClient := client.NewClientOrDie()
	dao, err := NewDaoBuilds()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	handler := NewBuildHandler(k8sClient, dao)

	for _, c := range cases {
		var reader io.Reader
		if len(c.body) > 0 {
			reader = strings.NewReader(c.body)
		}

		req, err := http.NewRequest("POST", "http://localhost"+c.url, reader)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept-Encoding", "gzip, deflate, sdch")

		w := httptest.NewRecorder()

		handler(w, req)

		if w.Code != c.expect {
			t.Logf("Expected HTTP Status Code of %d.  Got %d.  URL = %s\n", c.expect, w.Code, c.url)
			failed = true
		}

		if failed {
			t.Fail()
		}
	}

}
