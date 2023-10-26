NAME = zk-calc-app
IMAGE_NAME = zk-calc-app
IMAGE_VERSION = 0.0.4-test

LOCATION ?= us-west1
PROJECT_ID ?= zerok-dev
REPOSITORY ?= test-apps-public

sync:
	go get -v ./...

build: sync
	go build -v -o $(NAME) main.go

docker-build: sync
	CGO_ENABLED=0 GOOS=linux $(ARCH) go build -v -o $(NAME) main.go
	docker build --no-cache -t $(IMAGE_PREFIX)$(IMAGE_NAME):$(IMAGE_VERSION) .

docker-build-gke: IMAGE_PREFIX := $(LOCATION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/
docker-build-gke: ARCH := GOARCH=amd64
docker-build-gke: docker-build

docker-push-gke: IMAGE_PREFIX := $(LOCATION)-docker.pkg.dev/$(PROJECT_ID)/$(REPOSITORY)/
docker-push-gke:
	docker push $(IMAGE_PREFIX)$(IMAGE_NAME):$(IMAGE_VERSION)

docker-build-push-gke: docker-build-gke docker-push-gke


#Proto generation
generate-proto:
	protoc --go_out=. --go-grpc_out=. calculator.proto

#Testing
grpcurl-list:
	grpcurl -plaintext localhost:50051 list

grpcurl-list-methods:
	grpcurl -plaintext localhost:50051 list calculator.Calculator

grpcurl-add:
	grpcurl -plaintext -d '{"a": 5, "b": 3}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.Add

grpcurl-other-add:
	grpcurl -plaintext -d '{"a": 5, "b": 3, "forward": "50052"}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.CallOtherAdd

grpcurl-subtract:
	grpcurl -plaintext -d '{"a": 5, "b": 3}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.Subtract

grpcurl-other-subtract:
	grpcurl -plaintext -d '{"a": 5, "b": 3, "forward": "50052"}' -H "traceparent: 00-0018317d2edc7d4cfa75b39e88c2d570-176069a07e84acaf-01" localhost:50051 calculator.Calculator.CallOtherSubtract
