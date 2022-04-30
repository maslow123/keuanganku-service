proto:
	# protoc --proto_path=pkg pkg/pb/*.proto --go_out=plugins=grpc:.
	protoc \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/gogoproto \
		-I pkg --go_out=plugins=grpc:. pkg/pb/*.proto

pos-proto:
	# protoc pkg/pb/*.proto --go_out=plugins=grpc:.
	protoc \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2 \
		-I D:/me/exercise/go/pkg/mod/github.com/gogo/protobuf@v1.3.2/gogoproto \
		-I. --go_out=plugins=grpc:. pkg/pb/*.proto

buildapi:
	docker-compose build --no-cache api

runapi: buildapi
	docker-compose up -d api

server:
	go run cmd/main.go

test:
	go test -v ./... -coverprofile cover.out