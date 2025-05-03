.PHONY: tidy
tidy:
	@go mod tidy -v

.PHONY: docker
docker:
	@rm core-engine || true
	@docker rmi -f gocopilot/core-engine:v0.0.1
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o core-engine .
	@#GOOS=linux GOARCH=$(uname -m) go build -tags=k8s -o core-engine .
	@docker build -t gocopilot/core-engine:v0.0.1 .

.PHONY: grpc
grpc:
	@npx buf generate api/proto