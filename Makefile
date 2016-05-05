
.PHONY: all build restore

all: container

container: build
	godep go test
	godep go test ./pkg/...


build: restore
	##CGO_ENABLED=0 godep go build -a -ldflags '-s' (can't use this 'cause we have a non-portable ref to 'user')
	godep go vet
	golint
	godep go build
	cp jeeves docker

more_cool_tests:
	go test -bench-. -cpu=1,4,16 -benchmem

escape_analysis:
	go test -gcflags=-m -bench=something
	# see https://golang.org/pkg/net/http/pprof/
	go tool pprof http://localhost:9090/debug/pprof/profile
	go tool pprof http://localhost:9090/debug/pprof/heap

deploy:
	go test -tags=integration

restore: Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt

Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt:
	godep restore
