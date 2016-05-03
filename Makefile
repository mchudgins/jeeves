
.PHONY: all build restore

all: build

build: restore
	#CGO_ENABLED=0 godep go build -a -ldflags '-s' (can't use this 'cause we have a non-portable ref to 'user')
	godep go build
	cp jeeves docker

restore: Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt

Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt:
	godep restore
