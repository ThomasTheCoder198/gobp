# gobp

Opinionated Go project boilerplate generator. Run `gobp new` and answer a few
questions — you get a buildable, structured Go service skeleton in seconds.

## Install

**Via `go install` (recommended):**

```bash
go install github.com/tienanhnguyen999/gobp/cmd/gobp@latest
```

**From a release binary:** download the pre-built binary for your platform from
the [Releases](https://github.com/tienanhnguyen999/gobp/releases) page and place
it on your `$PATH`.

## Quickstart

### Interactive wizard

```bash
gobp new
```

The wizard walks you through 9 steps (name, module, framework, databases, SDKs,
patterns, addons, WebSocket, confirm). At any step press `Esc` to go back and
change a previous choice — your selections are preserved.

### Non-interactive (flags only)

```bash
gobp new --name myapi --module github.com/me/myapi --framework gin
cd myapi
go build ./...
go run ./cmd/server     # → http://localhost:8080/healthz
```

## Flags

| Flag | Default | Notes |
|---|---|---|
| `--name` | — | Project directory name; omit to launch the wizard |
| `--module` | `github.com/example/<name>` | Go module path |
| `--framework` | `gin` | `gin`, `echo`, or `fiber` |
| `--db` | — | Multi-value: `postgres`, `mysql`, `sqlite`, `mongo`, `redis`, `cassandra` |
| `--sdk` | — | Multi-value: `openai`, `stripe` |
| `--addon` | — | Multi-value: `docker`, `compose` (auto-added with any DB), `githubactions` |
| `--pattern` | — | Multi-value: `worker` |
| `--websocket` | `false` | Scaffold a WebSocket handler |
| `--git` | `init` | `init`, `commit`, or `none` |
| `--no-tidy` | `false` | Skip `go mod tidy` after generation |
| `--dry-run` | `false` | Print the render plan without writing files |
| `--out` | `./<name>` | Override the output directory |

## Introspect the registry

```bash
gobp list frameworks    # gin, echo, fiber
gobp list dbs           # postgres, mysql, sqlite, mongo, redis, cassandra
gobp list sdks          # openai, stripe
gobp list addons        # docker, compose, githubactions
gobp list patterns      # worker
gobp list features      # websocket
```

## Generated layout

```
<name>/
├── cmd/
│   ├── server/         # HTTP entry point (wire.go + wire_gen.go)
│   ├── cli/            # CLI entry point
│   └── worker/         # Background worker (only with --pattern worker)
├── internal/server/
│   ├── controller/     # Route handlers
│   ├── service/        # Business logic interfaces
│   ├── repository/     # Data access interfaces
│   ├── infrastructure/ # SDK clients (OpenAI, Stripe …)
│   ├── middleware/     # Trace-ID, logger, error, recovery
│   ├── handler/        # WebSocket handler (only with --websocket)
│   ├── dto/ model/ entity/
├── pkg/
│   ├── config/         # Viper-typed config + config.yaml
│   ├── database/<db>/  # One connection stub per selected DB
│   └── queue/          # consumer.go + producer.go (with --pattern worker)
├── errors/             # AppError + stable numeric error codes
├── log/                # Zap logger with trace_id context helpers
├── deploy/             # Dockerfile + docker-compose.yml
├── .github/workflows/  # ci.yml (with --addon githubactions)
├── mock/ migration/ schema/  # empty starter folders
└── Makefile README.md .gitignore
```

## Develop

```bash
make build       # build ./bin/gobp (version stamped from git tag)
make test        # run all tests
make fmt         # gofmt -s -w .
make lint        # go vet ./...
make install     # go install with version stamp
```

## Release a new version

1. Tag the commit:
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```
2. Build the release binary:
   ```bash
   make build TAG=v0.2.0
   # or: go build -ldflags "-X main.Version=v0.2.0" -o bin/gobp ./cmd/gobp
   ```
3. The binary reports its version via `gobp --version`.

For automated cross-platform builds, wire up
[GoReleaser](https://goreleaser.com) to the `release` Makefile target.

## License

MIT. See [LICENSE](LICENSE).
