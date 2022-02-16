test: lint
	@go test -count=1 -race ./...

build-server:
	@go build -o testify-usage-example main.go

run-server: build-server
	@go run main.go

lint: ## Lint the files
	@golangci-lint run
