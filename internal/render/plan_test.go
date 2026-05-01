package render

import (
	"testing"

	"github.com/tienanhnguyen999/gobp/internal/registry"
	"github.com/tienanhnguyen999/gobp/internal/selection"
)

func TestBuildPlan_Minimal(t *testing.T) {
	reg, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	sel := selection.Selection{
		Name:      "x",
		Module:    "github.com/example/x",
		Framework: "gin",
		GoVersion: "1.25",
	}
	p, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"go.mod", "Makefile", ".env", ".env.example", "README.md", ".gitignore",
		"cmd/server/main.go", "cmd/cli/main.go",
		"cmd/server/wire.go", "cmd/server/wire_gen.go",
		"pkg/config/config.go", "config.yaml",
		"errors/errors.go", "errors/codes.go",
		"log/logger.go",
		"internal/server/middleware/trace_id_middleware.go",
		"internal/server/controller/health_controller.go",
	}
	for _, w := range want {
		found := false
		for _, f := range p.Files {
			if f.Target == w {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing %q in plan", w)
		}
	}
}

func TestBuildPlan_AddsDBs(t *testing.T) {
	reg, _ := registry.Load()
	sel := selection.Selection{
		Name:      "x",
		Module:    "github.com/example/x",
		Framework: "gin",
		GoVersion: "1.25",
		DBs:       []string{"postgres", "redis"},
	}
	p, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}
	for _, w := range []string{"pkg/database/postgres/postgres.go", "pkg/database/redis/redis.go", "deploy/docker-compose.yml"} {
		found := false
		for _, f := range p.Files {
			if f.Target == w {
				found = true
			}
		}
		if !found {
			t.Errorf("missing %q in plan", w)
		}
	}
}

func TestBuildPlan_WorkerPattern(t *testing.T) {
	reg, _ := registry.Load()
	sel := selection.Selection{
		Name: "x", Module: "github.com/example/x", Framework: "gin", GoVersion: "1.25",
		Patterns: []string{"worker"},
	}
	p, err := BuildPlan(sel, reg)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"cmd/worker/main.go", "pkg/queue/consumer.go", "pkg/queue/producer.go"}
	for _, w := range want {
		found := false
		for _, f := range p.Files {
			if f.Target == w {
				found = true
			}
		}
		if !found {
			t.Errorf("missing %q in plan", w)
		}
	}
}

func TestBuildPlan_NoWorkerPatternByDefault(t *testing.T) {
	reg, _ := registry.Load()
	sel := selection.Selection{Name: "x", Module: "github.com/example/x", Framework: "gin", GoVersion: "1.25"}
	p, _ := BuildPlan(sel, reg)
	for _, f := range p.Files {
		if f.Target == "cmd/worker/main.go" {
			t.Errorf("worker should not be in plan when pattern not selected")
		}
	}
}

func TestBuildPlan_FrameworkTemplatePaths(t *testing.T) {
	reg, _ := registry.Load()
	for _, fw := range []string{"gin", "echo", "fiber"} {
		t.Run(fw, func(t *testing.T) {
			sel := selection.Selection{
				Name: "x", Module: "github.com/example/x", Framework: fw, GoVersion: "1.25",
			}
			p, err := BuildPlan(sel, reg)
			if err != nil {
				t.Fatal(err)
			}
			fwTargets := map[string]bool{
				"cmd/server/main.go":                                    true,
				"internal/server/middleware/trace_id_middleware.go":      true,
				"internal/server/middleware/logger_middleware.go":        true,
				"internal/server/middleware/error_middleware.go":         true,
				"internal/server/middleware/recovery_middleware.go":      true,
				"internal/server/controller/health_controller.go":       true,
			}
			for _, f := range p.Files {
				if fwTargets[f.Target] {
					wantPrefix := fw + "/"
					if len(f.Template) < len(wantPrefix) || f.Template[:len(wantPrefix)] != wantPrefix {
						t.Errorf("target %q: template %q should start with %q", f.Target, f.Template, wantPrefix)
					}
				}
			}
		})
	}
}
