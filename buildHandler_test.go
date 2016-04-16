package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPost(t *testing.T) {
	cases := []struct {
		url    string
		body   string
		expect int
	}{
		{url: "/builds/golang:latest", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar/gorf", body: "", expect: 400},
		{url: "/builds/golang:latest/fubar?gorf", body: "", expect: 400},
	}

	var failed bool

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

		BuildHandler(w, req)

		if w.Code != c.expect {
			t.Logf("Expected HTTP Status Code of %d.  Got %d.  URL = %s\n", c.expect, w.Code, c.url)
			failed = true
		}

		if failed {
			t.Fail()
		}
	}

}
