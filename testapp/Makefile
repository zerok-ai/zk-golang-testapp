
sync:
	go get -v ./...

build: sync
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o testapp ./main.go

docker-build:
	docker build -t us-west1-docker.pkg.dev/zerok-dev/golang-testapp/testapp:latest ./