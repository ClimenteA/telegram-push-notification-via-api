build:
	GOOS=linux GOARCH=amd64 go build -o dist/server server.go
	docker-compose build --no-cache

dev:
	docker-compose up

run: 
	docker-compose up -d