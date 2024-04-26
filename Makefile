.PHONY: bench build build-linux install lint run race test update

git_sha := $(shell git rev-parse --short HEAD)
ld_flags := -ldflags "-X 'main.GitSHA=$(git_sha)'"
build_cmd := go build $(ld_flags) -o bin/jestclout ./cmd/jestclout

# Build

bench:
	go test -bench=. -cpu=1,4,16 -benchmem ./...

build:
	$(build_cmd)

build-linux:
	CGO_ENABLED=0 GOOS=linux $(build_cmd)

install:
	go install ./...

lint:
	golangci-lint run

run:
	go run $(ld_flags) ./cmd/jestclout

race:
	go test -race ./...

test:
	go test -cover ./...

update:
	go get -u all
