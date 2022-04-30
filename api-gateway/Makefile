proto:
	protoc --proto_path=pkg pkg/**/pb/*.proto --go_out=plugins=grpc:.
	protoc \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/gogoproto \
		-I pkg pkg/transactions/pb/*.proto --go_out=plugins=grpc:.

buildapigateway:
	docker-compose build api-gateway

runapigateway: buildapigateway
	docker-compose up -d api-gateway

server:
	go run cmd/main.go

test:
	go test -v ./... -coverprofile cover.out