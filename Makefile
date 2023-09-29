
build:
	cd testapp && make build

buildAndPush: build
	cd testapp && make docker-build
	docker push us-west1-docker.pkg.dev/zerok-dev/golang-testapp/testapp:latest

deploy:
	kubectl apply -k ./k8s