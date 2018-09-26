lint:
	docker run --rm -t \
	-v $(GOPATH)/src/github.com/techdroplabs/cloudwatch-cleaner:/go/src/github.com/techdroplabs/cloudwatch-cleaner \
	-w /go/src/github.com/techdroplabs/cloudwatch-cleaner \
	golangci/golangci-lint run
.PHONY: lint

test:
	docker run --rm -t \
	-v $(GOPATH)/src/github.com/techdroplabs/cloudwatch-cleaner:/go/src/github.com/techdroplabs/cloudwatch-cleaner \
	-w /go/src/github.com/techdroplabs/cloudwatch-cleaner \
	golang:latest go test -timeout 20s -race -coverprofile coverage.txt -covermode=atomic ./...
.PHONY: test

build:
	GOOS=linux go build main.go
	zip handler.zip ./lambda_handler
