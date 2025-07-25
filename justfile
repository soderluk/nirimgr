version := `git describe --tags --always`
commit := `git rev-parse --short HEAD`
buildDate := `date -u '+%Y-%m-%dT%H:%M:%SZ'`
ldflags := "-s -w -X github.com/soderluk/nirimgr/config.Version=" + version + " -X github.com/soderluk/nirimgr/config.CommitSHA=" + commit + " -X github.com/soderluk/nirimgr/config.BuildDate=" + buildDate

help:
    @just --list

version:
    @echo {{ version }}

fmt:
    go fmt ./...

vet:
    go vet ./...

build: fmt vet
    @go build -ldflags "{{ ldflags }}" .

install:
    @go install -ldflags "{{ ldflags }}"

test:
    @go test ./... -coverprofile cover.out

run RUNARGS:
    @go run -ldflags "{{ ldflags }}" ./main.go {{ RUNARGS }}

coverage:
    @go tool cover -html=cover.out
