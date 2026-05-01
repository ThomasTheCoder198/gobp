// Package render builds and executes the render plan that turns a Selection
// into files on disk.
package render

import (
	"fmt"
	"sort"

	"github.com/tienanhnguyen999/gobp/internal/registry"
	"github.com/tienanhnguyen999/gobp/internal/selection"
	"github.com/tienanhnguyen999/gobp/internal/util"
)

// Plan is the resolved render plan: an ordered list of file targets, each
// produced by either a single template or a list of section fragments.
type Plan struct {
	Selection selection.Selection
	Files     []FileEntry
}

// FileEntry describes one output file.
type FileEntry struct {
	// Target is the path relative to the project root.
	Target string

	// Template is the embedded template path (under internal/templates/).
	// Either Template or Static must be set, not both.
	Template string

	// Static is verbatim file content (used for things like .gitkeep markers).
	Static []byte

	// Format requests gofmt formatting on the rendered output. Set for *.go.
	Format bool

	// Executable marks the file as +x.
	Executable bool
}

// BuildPlan resolves Selection + Registry into a deterministic Plan.
//
// The plan is built by enumerating well-known project files and toggling each
// one based on the selection. This is intentionally simpler than a fully
// section-based composition engine: every conditional branch lives inside
// the template files, not in the plan builder. To add a brand-new top-level
// file, add an entry here.
func BuildPlan(sel selection.Selection, reg *registry.Registry) (*Plan, error) {
	p := &Plan{Selection: sel}

	fw := sel.Framework

	// --- Always-present project files ---
	add := func(target, tmpl string, format bool) {
		p.Files = append(p.Files, FileEntry{Target: target, Template: tmpl, Format: format})
	}
	fwAdd := func(target, tmpl string, format bool) {
		add(target, fw+"/"+tmpl, format)
	}

	add("go.mod", "go.mod.tmpl", false)
	add("README.md", "README.md.tmpl", false)
	add("Makefile", "Makefile.tmpl", false)
	add(".gitignore", ".gitignore.tmpl", false)
	add(".env.example", ".env.example.tmpl", false)
	add(".env", ".env.tmpl", false)

	// config.yaml lives at the project root so it is easy to find and edit.
	add("config.yaml", "pkg/config/config.yaml.tmpl", false)

	// Entry points
	fwAdd("cmd/server/main.go", "cmd/server/main.go.tmpl", true)
	fwAdd("cmd/server/server.go", "cmd/server/server.go.tmpl", true)
	add("cmd/server/README.md", "cmd/server/README.md.tmpl", false)
	add("cmd/cli/main.go", "cmd/cli/main.go.tmpl", true)
	add("cmd/cli/README.md", "cmd/cli/README.md.tmpl", false)

	// Config
	add("pkg/config/config.go", "pkg/config/config.go.tmpl", true)
	add("pkg/config/README.md", "pkg/config/README.md.tmpl", false)

	// Errors
	add("errors/errors.go", "errors/errors.go.tmpl", true)
	add("errors/codes.go", "errors/codes.go.tmpl", true)
	add("errors/README.md", "errors/README.md.tmpl", false)

	// Logger
	add("log/logger.go", "log/logger.go.tmpl", true)
	add("log/README.md", "log/README.md.tmpl", false)

	// Utilities
	add("pkg/utils/string.go", "pkg/utils/string.go.tmpl", true)
	add("pkg/utils/pointer.go", "pkg/utils/pointer.go.tmpl", true)
	add("pkg/utils/slice.go", "pkg/utils/slice.go.tmpl", true)
	add("pkg/utils/time.go", "pkg/utils/time.go.tmpl", true)
	add("pkg/utils/pagination.go", "pkg/utils/pagination.go.tmpl", true)
	add("pkg/utils/response.go", "pkg/utils/response.go.tmpl", true)
	add("pkg/utils/request.go", "pkg/utils/request.go.tmpl", true)
	add("pkg/utils/validation.go", "pkg/utils/validation.go.tmpl", true)
	add("pkg/utils/retry.go", "pkg/utils/retry.go.tmpl", true)
	add("pkg/utils/convert.go", "pkg/utils/convert.go.tmpl", true)
	add("pkg/utils/README.md", "pkg/utils/README.md.tmpl", false)

	// Server middleware (CORS always included)
	fwAdd("internal/server/middleware/cors_middleware.go", "internal/server/middleware/cors_middleware.go.tmpl", true)
	fwAdd("internal/server/middleware/trace_id_middleware.go", "internal/server/middleware/trace_id_middleware.go.tmpl", true)
	fwAdd("internal/server/middleware/logger_middleware.go", "internal/server/middleware/logger_middleware.go.tmpl", true)
	fwAdd("internal/server/middleware/error_middleware.go", "internal/server/middleware/error_middleware.go.tmpl", true)
	fwAdd("internal/server/middleware/recovery_middleware.go", "internal/server/middleware/recovery_middleware.go.tmpl", true)
	add("internal/server/middleware/README.md", "internal/server/middleware/README.md.tmpl", false)

	// Health controller
	fwAdd("internal/server/controller/health_controller.go", "internal/server/controller/health_controller.go.tmpl", true)
	add("internal/server/controller/README.md", "internal/server/controller/README.md.tmpl", false)

	// Wire DI — colocated with cmd/server/main.go so they share package main.
	add("cmd/server/wire.go", "wire.go.tmpl", true)
	add("cmd/server/wire_gen.go", "wire_gen.go.tmpl", true)

	// --- Empty-folder README markers ---
	for _, dir := range []string{
		"internal/server/service",
		"internal/server/repository",
		"internal/server/infrastructure",
		"internal/server/dto",
		"internal/server/model",
		"internal/server/entity",
		"mock",
		"schema",
	} {
		add(dir+"/README.md", dir+"/README.md.tmpl", false)
	}

	// Migration folder: example SQL file + README
	add("migration/README.md", "migration/README.md.tmpl", false)
	add("migration/0001_init.sql", "migration/0001_init.sql.tmpl", false)

	// --- Conditional: WebSocket handler ---
	if sel.Websocket {
		fwAdd("internal/server/handler/ws_handler.go", "internal/server/handler/ws_handler.go.tmpl", true)
		add("internal/server/handler/README.md", "internal/server/handler/README.md.tmpl", false)
	}

	// --- Conditional: per-DB stubs ---
	for _, db := range sel.DBs {
		switch db {
		case "postgres":
			add("pkg/database/postgres/postgres.go", "pkg/database/postgres/postgres.go.tmpl", true)
			add("pkg/database/postgres/README.md", "pkg/database/postgres/README.md.tmpl", false)
		case "redis":
			add("pkg/database/redis/redis.go", "pkg/database/redis/redis.go.tmpl", true)
			add("pkg/database/redis/README.md", "pkg/database/redis/README.md.tmpl", false)
		case "mysql":
			add("pkg/database/mysql/mysql.go", "pkg/database/mysql/mysql.go.tmpl", true)
			add("pkg/database/mysql/README.md", "pkg/database/mysql/README.md.tmpl", false)
		case "mongo":
			add("pkg/database/mongo/mongo.go", "pkg/database/mongo/mongo.go.tmpl", true)
			add("pkg/database/mongo/README.md", "pkg/database/mongo/README.md.tmpl", false)
		case "sqlite":
			add("pkg/database/sqlite/sqlite.go", "pkg/database/sqlite/sqlite.go.tmpl", true)
			add("pkg/database/sqlite/README.md", "pkg/database/sqlite/README.md.tmpl", false)
		case "cassandra":
			add("pkg/database/cassandra/cassandra.go", "pkg/database/cassandra/cassandra.go.tmpl", true)
			add("pkg/database/cassandra/README.md", "pkg/database/cassandra/README.md.tmpl", false)
		default:
			return nil, fmt.Errorf("unsupported db %q in plan", db)
		}
	}

	// --- Conditional: per-SDK stubs ---
	for _, sdk := range sel.SDKs {
		switch sdk {
		case "openai":
			add("internal/server/infrastructure/openai/openai.go", "internal/server/infrastructure/openai/openai.go.tmpl", true)
			add("internal/server/infrastructure/openai/README.md", "internal/server/infrastructure/openai/README.md.tmpl", false)
		case "stripe":
			add("internal/server/infrastructure/stripe/stripe.go", "internal/server/infrastructure/stripe/stripe.go.tmpl", true)
			add("internal/server/infrastructure/stripe/README.md", "internal/server/infrastructure/stripe/README.md.tmpl", false)
		default:
			return nil, fmt.Errorf("unsupported sdk %q in plan", sdk)
		}
	}

	// --- Conditional: worker pattern ---
	if util.StringIn(sel.Patterns, "worker") {
		add("cmd/worker/main.go", "cmd/worker/main.go.tmpl", true)
		add("cmd/worker/README.md", "cmd/worker/README.md.tmpl", false)
		add("pkg/queue/consumer.go", "pkg/queue/consumer.go.tmpl", true)
		add("pkg/queue/producer.go", "pkg/queue/producer.go.tmpl", true)
		add("pkg/queue/README.md", "pkg/queue/README.md.tmpl", false)
	}

	// --- Conditional: docker addon ---
	if util.StringIn(sel.Addons, "docker") {
		add("deploy/Dockerfile", "deploy/Dockerfile.tmpl", false)
	}
	// Compose: auto-on when any DB selected, or explicitly via --addon compose.
	if util.StringIn(sel.Addons, "compose") || len(sel.DBs) > 0 {
		add("deploy/docker-compose.yml", "deploy/docker-compose.yml.tmpl", false)
	}
	// deploy/README.md: whenever the deploy/ directory will exist.
	if util.StringIn(sel.Addons, "docker") || util.StringIn(sel.Addons, "compose") || len(sel.DBs) > 0 {
		add("deploy/README.md", "deploy/README.md.tmpl", false)
	}
	// GitHub Actions
	if util.StringIn(sel.Addons, "githubactions") {
		add(".github/workflows/ci.yml", ".github/workflows/ci.yml.tmpl", false)
	}

	// Dedupe by target (later entries win — should never happen, but be safe).
	sort.SliceStable(p.Files, func(i, j int) bool { return p.Files[i].Target < p.Files[j].Target })
	out := make([]FileEntry, 0, len(p.Files))
	seen := map[string]int{}
	for _, f := range p.Files {
		if idx, ok := seen[f.Target]; ok {
			out[idx] = f
			continue
		}
		seen[f.Target] = len(out)
		out = append(out, f)
	}
	p.Files = out

	_ = reg // reserved for richer manifest-driven composition
	return p, nil
}

// PrintPlan writes a human-readable summary of the plan.
func PrintPlan(w interface{ Write([]byte) (int, error) }, p *Plan) error {
	for _, f := range p.Files {
		line := f.Target + "\n"
		if _, err := w.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
