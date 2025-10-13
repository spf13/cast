[private]
default:
    just --list

# run all checks
check: test lint

# run tests
test:
    go test -shuffle on -race -v ./...

# run linter
lint:
    golangci-lint run

# format code
fmt:
    golangci-lint fmt
