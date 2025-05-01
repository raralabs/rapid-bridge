swagger:
	@swag init -g ./internal/route/swagger.go --parseDependency --parseInternal -d ./

server:
	@go run cmd/server/server.go

build:
	@go build -o rapid-bridge cmd/main.go

watch-server:
	@wgo run cmd/server/server.go

gotidy:
	@GOPRIVATE=github.com/raralabs go mod tidy

proto: generate/proto inject-tag
	echo "Generated"

generate/proto:
	go run github.com/bufbuild/buf/cmd/buf@v1.32.1 generate --exclude-path grpc/proto/google
	
inject-tag:
	protoc-go-inject-tag -input="grpc/pb/*.pb.go"

