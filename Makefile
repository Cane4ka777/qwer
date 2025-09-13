APP_NAME=qwer-api
PORT?=8080

.PHONY: run build docker-build docker-run fmt tidy

run:
	go run main.go

build:
	go build -o bin/$(APP_NAME) main.go

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run --rm -p $(PORT):8080 $(APP_NAME):latest

fmt:
	gofmt -s -w .

tidy:
	go mod tidy
