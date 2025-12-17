version := `jj log --no-graph -r @ -T 'parents.map(|c| c.tags())'`
commit := `jj log -T 'commit_id.short() ++ "\n"' --no-graph | head -n1`
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
