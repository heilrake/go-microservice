PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .

.PHONY: generate-proto proto-docs swagger-docs

generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

# Generate HTML documentation from proto files.
# Requires Docker. Output: docs/proto-docs/index.html
proto-docs:
	docker run --rm \
		--platform linux/amd64 \
		-v $(PWD)/proto:/proto \
		-v $(PWD)/docs/proto-docs:/out \
		pseudomuto/protoc-gen-doc \
		--doc_opt=html,index.html \
		driver.proto user.proto trip.proto

# Regenerate Swagger docs for the API Gateway.
# Requires swag CLI: go install github.com/swaggo/swag/cmd/swag@latest
swagger-docs:
	cd services/api-gateway && swag init \
		-g main.go \
		-d .,../../shared/contracts,../../shared/types \
		--parseInternal \
		-o ./docs