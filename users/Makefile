proto:
	protoc pkg/pb/*.proto --go_out=plugins=grpc:.

buildapi:
	docker-compose build --no-cache api

runapi: 
	docker-compose up -d api

server:
	go run cmd/main.go

test:
	go test -v ./... -coverprofile cover.out