default:
    just --list

# Run all tests
test:
    go test ./...

# Run static checks available in the Go toolchain
lint:
    go vet ./...

# Format Go sources
format:
    gofmt -w $(find . -name '*.go' -not -path './vendor/*')

# Tidy module dependencies
tidy:
    go mod tidy

# Build the CLI into ./bin/sear
build:
    mkdir -p bin
    go build -o bin/sear ./cmd/sear

# Install the CLI into ~/.local/bin/sear
install:
    mkdir -p "$HOME/.local/bin"
    go build -o "$HOME/.local/bin/sear" ./cmd/sear

# Run the CLI, pass args after --, e.g. just run -- -q 'what is searxng'
run *ARGS:
    go run ./cmd/sear {{ARGS}}
