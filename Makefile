build:
	@go build -o bin/cmd cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/cmd