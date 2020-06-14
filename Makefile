.PHONY: build
build:
	go mod download && go mod vendor && go build -mod=vendor -o ./bin/redditclone ./cmd/redditclone

.PHONY: run
run:
	go run cmd/redditclone/main.go

.PHONY: up
up:
	make build && make run

.PHONY: test
test:
	go test -v -mod=vendor -coverpkg=./... ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html

.PHONY: docker-up
docker-up:
	docker-compose up -d --build --remove-orphans

.PHONY: docker-down
docker-down:
	docker-compose down