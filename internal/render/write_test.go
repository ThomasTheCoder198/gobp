package render

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tienanhnguyen999/gobp/internal/registry"
	"github.com/tienanhnguyen999/gobp/internal/selection"
)

// TestExecute_MinimalRendersAndFormats renders the minimal gin project to a
// tempdir and verifies that the expected files exist. Per-file gofmt formatting
// is exercised inside Execute itself; if any .go template were malformed it
// would fail here with a clear error.
func TestExecute_MinimalRendersAndFormats(t *testing.T) {
	reg, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	sel := selection.Selection{
		Name:      "x",
		Module:    "github.com/example/x",
		Framework: "gin",
		GoVersion: "1.25",
		Year:      2026,
	}
	plan, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}

	tmp, err := os.MkdirTemp("", "gobp-render-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	target := filepath.Join(tmp, "out")

	if err := Execute(plan, target); err != nil {
		t.Fatalf("execute: %v", err)
	}

	required := []string{
		"go.mod", "Makefile", ".env", ".env.example", "README.md", ".gitignore",
		"cmd/server/main.go", "cmd/server/wire.go", "cmd/server/wire_gen.go",
		"cmd/cli/main.go",
		"pkg/config/config.go", "config.yaml",
		"errors/errors.go", "errors/codes.go",
		"log/logger.go",
		"internal/server/middleware/trace_id_middleware.go",
		"internal/server/middleware/logger_middleware.go",
		"internal/server/middleware/recovery_middleware.go",
		"internal/server/middleware/error_middleware.go",
		"internal/server/controller/health_controller.go",
		"mock/README.md", "migration/README.md", "schema/README.md",
	}
	for _, p := range required {
		path := filepath.Join(target, p)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
}

func TestExecute_RefusesNonEmptyTarget(t *testing.T) {
	reg, _ := registry.Load()
	sel := selection.Selection{Name: "x", Module: "github.com/example/x", Framework: "gin", GoVersion: "1.25"}
	plan, _ := BuildPlan(sel, reg)

	tmp, _ := os.MkdirTemp("", "gobp-busy-")
	defer func() { _ = os.RemoveAll(tmp) }()
	if err := os.WriteFile(filepath.Join(tmp, "preexisting"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Execute(plan, tmp); err == nil {
		t.Errorf("expected error when target dir is non-empty")
	}
}

func TestExecute_EchoRendersAndFormats(t *testing.T) {
	reg, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	sel := selection.Selection{
		Name:      "x",
		Module:    "github.com/example/x",
		Framework: "echo",
		GoVersion: "1.25",
		Year:      2026,
	}
	plan, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}

	tmp, err := os.MkdirTemp("", "gobp-echo-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	target := filepath.Join(tmp, "out")

	if err := Execute(plan, target); err != nil {
		t.Fatalf("execute echo: %v", err)
	}

	required := []string{
		"go.mod", "cmd/server/main.go",
		"internal/server/middleware/trace_id_middleware.go",
		"internal/server/middleware/logger_middleware.go",
		"internal/server/middleware/recovery_middleware.go",
		"internal/server/middleware/error_middleware.go",
		"internal/server/controller/health_controller.go",
	}
	for _, p := range required {
		path := filepath.Join(target, p)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
}

func TestExecute_FiberRendersAndFormats(t *testing.T) {
	reg, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	sel := selection.Selection{
		Name:      "x",
		Module:    "github.com/example/x",
		Framework: "fiber",
		GoVersion: "1.25",
		Year:      2026,
	}
	plan, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}

	tmp, err := os.MkdirTemp("", "gobp-fiber-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	target := filepath.Join(tmp, "out")

	if err := Execute(plan, target); err != nil {
		t.Fatalf("execute fiber: %v", err)
	}

	required := []string{
		"go.mod", "cmd/server/main.go",
		"internal/server/middleware/trace_id_middleware.go",
		"internal/server/middleware/logger_middleware.go",
		"internal/server/middleware/recovery_middleware.go",
		"internal/server/middleware/error_middleware.go",
		"internal/server/controller/health_controller.go",
	}
	for _, p := range required {
		path := filepath.Join(target, p)
		if _, err := os.Stat(path); err != nil {
			t.Errorf("missing %s: %v", p, err)
		}
	}
}
