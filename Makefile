test:
	go test -coverprofile=coverage.out ./...

	go tool cover -func=coverage.out

build:
	docker-compose build

up:
	docker-compose up -d