
.PHONY: all build restore

all: build

build: restore
	godep go build

restore: Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt

Godeps/_workspace/src/github.com/aws/aws-sdk-go/NOTICE.txt:	
	godep restore
