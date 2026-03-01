# Build the binary
build:
	@go build -o ./bin/vinomclient ./main.go

# Run tests
test:
	go test -v ./...

# Run the application
run: build
	@./bin/vinomclient
