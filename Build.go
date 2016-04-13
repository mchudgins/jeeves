package main

import "time"

type Build struct {
	Build       string
	LastUpdated time.Time
	Namespace   string
	Number      int
}
