.PHONY: docker
docker:
	@rm testengine || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o testengine .
	@docker rmi -f testcopilot/testengine:v0.0.1
	@docker build -t testcopilot/testengine-live:v0.0.1 .

.PHONY: grpc
grpc:
	@npx buf generate api/proto