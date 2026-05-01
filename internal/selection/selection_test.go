package selection

import (
	"strings"
	"testing"

	"github.com/tienanhnguyen999/gobp/internal/registry"
)

func mustRegistry(t *testing.T) *registry.Registry {
	t.Helper()
	r, err := registry.Load()
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}
	return r
}

func TestResolve_AppliesDefaults(t *testing.T) {
	sel := Selection{Name: "myapp"}
	if err := Resolve(&sel, nil); err != nil {
		t.Fatal(err)
	}
	if sel.Framework != "gin" {
		t.Errorf("framework: want gin, got %q", sel.Framework)
	}
	if sel.GoVersion != "1.25" {
		t.Errorf("goversion: want 1.25, got %q", sel.GoVersion)
	}
	if sel.Git != "init" {
		t.Errorf("git: want init, got %q", sel.Git)
	}
	if sel.Module == "" {
		t.Errorf("module should default when name is set")
	}
}

func TestResolve_AutoCompose(t *testing.T) {
	sel := Selection{Name: "x", DBs: []string{"postgres"}}
	if err := Resolve(&sel, nil); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, a := range sel.Addons {
		if a == "compose" {
			found = true
		}
	}
	if !found {
		t.Errorf("compose addon should auto-add when DB selected; got %v", sel.Addons)
	}
}

func TestResolve_DedupesSlices(t *testing.T) {
	sel := Selection{Name: "x", DBs: []string{"postgres", "postgres", "redis"}}
	if err := Resolve(&sel, nil); err != nil {
		t.Fatal(err)
	}
	if len(sel.DBs) != 2 {
		t.Errorf("dedup failed: %v", sel.DBs)
	}
}

func TestValidate_RejectsBadName(t *testing.T) {
	reg := mustRegistry(t)
	cases := []string{"", "with space", "9starts", "../etc", "a/b"}
	for _, n := range cases {
		sel := Selection{Name: n, Module: "github.com/x/y"}
		_ = Resolve(&sel, reg)
		err := Validate(&sel, reg)
		if err == nil {
			t.Errorf("expected error for name %q", n)
		}
	}
}

func TestValidate_RejectsBadModule(t *testing.T) {
	reg := mustRegistry(t)
	sel := Selection{Name: "ok", Module: "not a path"}
	_ = Resolve(&sel, reg)
	if err := Validate(&sel, reg); err == nil {
		t.Errorf("expected validation error for bad module path")
	}
}

func TestValidate_RejectsUnknownIDs(t *testing.T) {
	reg := mustRegistry(t)
	sel := Selection{Name: "ok", Module: "github.com/x/y", DBs: []string{"foobar"}}
	_ = Resolve(&sel, reg)
	err := Validate(&sel, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown db") {
		t.Errorf("expected unknown db error; got %v", err)
	}
}

func TestValidate_AcceptsKnownIDs(t *testing.T) {
	reg := mustRegistry(t)
	sel := Selection{
		Name:      "ok",
		Module:    "github.com/x/y",
		Framework: "gin",
		DBs:       []string{"postgres", "redis"},
		SDKs:      []string{"openai"},
		Patterns:  []string{"worker"},
	}
	if err := Resolve(&sel, reg); err != nil {
		t.Fatal(err)
	}
	if err := Validate(&sel, reg); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}
