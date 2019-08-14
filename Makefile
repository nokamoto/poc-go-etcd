.PHONY: docker

mod:
	rm -f go.mod go.sum
	go mod init github.com/nokamoto/poc-go-etcd

all: test bin docker

test: fmt
	go test .
	go mod tidy

fmt:
	gofmt -d -w .

bin:
	go build .

docker:
	docker-compose build

exec: all
	docker-compose down
	docker-compose up
