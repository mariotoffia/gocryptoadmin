test:
	@go test ./...
lint:
	@gocritic check ./...
dep:
	@GO111MODULE=on go get -v -u github.com/go-critic/go-critic/cmd/gocritic
	@go install github.com/segmentio/golines@latest