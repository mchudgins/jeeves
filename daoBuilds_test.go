package main

import (
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		in     Build
		expect string
	}{
		{in: Build{Build: "fubar", Namespace: "gorf", Number: 1}, expect: "first"},
		{in: Build{Build: "first", Namespace: "gorf", Number: 5}, expect: "first"},
		{in: Build{Build: "second", Namespace: "gorf", Number: 5}, expect: "first"},
		{in: Build{Build: "third", Namespace: "gorf", Number: 5}, expect: "first"},
		{in: Build{Build: "third", Namespace: "gorf", Number: 6}, expect: "first"},
		{in: Build{Build: "third", Namespace: "gorf", Number: 7}, expect: "first"},
		{in: Build{Build: "third", Namespace: "gorf", Number: 8}, expect: "first"},
	}

	dao, err := NewDaoBuilds()
	if err != nil {
		t.Errorf("unable to construct DaoBuild struct: %v", err)
	}

	for _, c := range cases {
		err = dao.Persist(&c.in)
		if err != nil {
			t.Errorf("unable to persist: %v", err)
		}
	}
}

func TestFetch(t *testing.T) {
	cases := []struct {
		in     Build
		expect string
	}{
		{in: Build{Build: "first"}, expect: "first"},
		{in: Build{Build: "second"}, expect: "second"},
		{in: Build{Build: "third"}, expect: "third"},
	}

	dao, err := NewDaoBuilds()
	if err != nil {
		t.Errorf("unable to construct DaoBuild struct: %v", err)
	}

	for _, c := range cases {
		obj, err := dao.Fetch("mch-dev0", c.in.Build, 1)
		if err != nil {
			t.Errorf("unable to Fetch %s: %v", c.in.Build, err)
		}
		if obj.Build != c.expect {
			t.Errorf("expected %s but got %s", c.expect, obj.Build)
		}
	}
}

func TestFetchByNamespace(t *testing.T) {
	dao, err := NewDaoBuilds()
	if err != nil {
		t.Errorf("unable to construct DaoBuild struct: %v", err)
	}

	builds, err := dao.FetchAllByNamespace("gorf")
	if err != nil {
		log.Fatal(err)
		t.Fail()
	}

	for _, b := range builds {
		log.Printf("build: %v", b)
	}
}
