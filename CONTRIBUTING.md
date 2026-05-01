# Contributing to gobp

## Prerequisites

| Tool | Minimum version | Install |
|------|----------------|---------|
| Go   | 1.25           | https://go.dev/dl |
| Git  | any recent     | https://git-scm.com |

No other tools are required to build or test the project.

## Clone and bootstrap

```bash
git clone https://github.com/tienanhnguyen999/gobp.git
cd gobp
go mod download        # fetch all dependencies
```

## Build

```bash
make build             # produces bin/gobp (or bin/gobp.exe on Windows)
```

Run it directly without installing:

```bash
./bin/gobp --help
./bin/gobp new --name myapi --framework gin --dry-run
```

## Install globally

```bash
make install           # go install ./cmd/gobp → adds gobp to your $GOPATH/bin
gobp --help
```

Make sure `$GOPATH/bin` (usually `~/go/bin`) is on your `PATH`.

## Run tests

```bash
make test              # go test ./...
```

All tests are unit tests with no external dependencies. They should pass on a
fresh clone with no extra setup.

## Code quality

```bash
make fmt               # format with gofmt -s -w .
make lint              # static analysis with go vet ./...
```

Run both before opening a pull request.

## Project layout

```
cmd/gobp/          CLI entry point
internal/
  cli/             Cobra commands (new, list)
  registry/        Manifest loader + component registry
  selection/       User selection model
  render/          Template engine + file writer
  postgen/         Post-generation hooks (go mod tidy, git init)
  wizard/          Bubble Tea interactive wizard (in progress)
internal/registry/manifests/   YAML manifests for each component
internal/render/templates/     Go text/template files for generated projects
```

## Trying a change end-to-end

1. Build: `make build`
2. Generate a test project: `./bin/gobp new --name testapp --framework gin --db postgres --out /tmp/testapp`
3. Verify the output compiles: `cd /tmp/testapp && go build ./...`
4. Clean up: `rm -rf /tmp/testapp`

## Common issues

**`gobp: command not found` after `make install`**
Add `$(go env GOPATH)/bin` to your `PATH`:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

**`go: go.mod requires go >= 1.25` error**
Update your Go toolchain: https://go.dev/dl

**Tests fail on Windows with path errors**
Run tests inside Git Bash or WSL — the test suite uses forward-slash paths.
