.PHONY: docker
docker:
	@rm demojob || true
	@GOOS=linux GOARCH=arm go build -o demojob .
	@docker rmi -f halfcoconut/hello-egg/demojob:v0.0.1
	@docker build -t halfcoconut/hello-egg/demojob-live:v0.0.1 .