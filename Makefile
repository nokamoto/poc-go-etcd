
all: test
	go build .

test: fmt
	go test .

fmt:
	gofmt -d -w .
