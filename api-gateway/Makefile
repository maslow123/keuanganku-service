proto:
	protoc --proto_path=pkg pkg/**/pb/*.proto --go_out=plugins=grpc:.
	protoc --proto_path=pkg pkg/transactions/pb/*.proto --go_out=plugins=grpc:.

buildapigateway:
	docker-compose build api-gateway

runapigateway: buildapigateway
	docker-compose up -d api-gateway

server:
	go run cmd/main.go

test:
	go test -v ./... -coverprofile cover.out