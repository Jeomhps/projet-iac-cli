BIN := "projet-iac-cli"

# Build into ./bin/
build:
	mkdir -p bin
	go build -o bin/{{BIN}} .

# Install to GOPATH/bin (or GOBIN)
install:
	go install .

# Run without building a binary first (pass args after --)
run *ARGS:
	go run . -- {{ARGS}}

# Format, tidy, test
fmt:
	go fmt ./...

tidy:
	go mod tidy

test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin
