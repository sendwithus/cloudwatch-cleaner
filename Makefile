lint:
	docker run --rm -t \
	-v $(GOPATH)/src/github.com/techdroplabs/cloudwatch-cleaner:/go/src/github.com/techdroplabs/cloudwatch-cleaner \
	-w /go/src/github.com/techdroplabs/cloudwatch-cleaner \
	golangci/golangci-lint run
.PHONY: lint

test:
	go test -timeout 20s -race -coverprofile coverage.txt -covermode=atomic ./...
.PHONY: test
