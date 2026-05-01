// Package selection holds the resolved user choices for a generated project.
package selection

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/mod/module"

	"github.com/tienanhnguyen999/gobp/internal/registry"
	"github.com/tienanhnguyen999/gobp/internal/util"
)

// Selection is the resolved set of user choices used by the renderer.
type Selection struct {
	Name      string
	Module    string
	Framework string
	DBs       []string
	SDKs      []string
	Addons    []string
	Patterns  []string
	Websocket bool

	Git    string // init, commit, none
	NoTidy bool
	DryRun bool

	GoVersion string
	Year      int
}

// Defaults returns a baseline Selection with safe defaults applied.
func Defaults() Selection {
	return Selection{
		Framework: "gin",
		Git:       "init",
		GoVersion: "1.25",
	}
}

var nameRE = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

// ValidateName checks that a project name is well-formed.
// Used by both Validate and the interactive wizard.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if !nameRE.MatchString(name) {
		return fmt.Errorf("must start with a letter and contain only letters, digits, _ or -")
	}
	if strings.ContainsAny(name, `/\`) || filepath.Base(name) != name {
		return fmt.Errorf("must not contain path separators")
	}
	return nil
}

// Resolve fills in missing fields with defaults and applies auto-rules.
//
// Auto-rules:
//   - if any DB is selected, the "compose" addon is auto-added
//   - duplicate entries in slice flags are de-duped
//   - module path defaults to github.com/<placeholder>/<name> if empty
func Resolve(sel *Selection, _ *registry.Registry) error {
	d := Defaults()

	if sel.Framework == "" {
		sel.Framework = d.Framework
	}
	if sel.Git == "" {
		sel.Git = d.Git
	}
	if sel.GoVersion == "" {
		sel.GoVersion = d.GoVersion
	}
	if sel.Year == 0 {
		sel.Year = 2026
	}

	sel.DBs = dedupSorted(sel.DBs)
	sel.SDKs = dedupSorted(sel.SDKs)
	sel.Addons = dedupSorted(sel.Addons)
	sel.Patterns = dedupSorted(sel.Patterns)

	// Auto-rule: any DB selected → ensure "compose" addon present.
	if len(sel.DBs) > 0 && !util.StringIn(sel.Addons, "compose") {
		sel.Addons = append(sel.Addons, "compose")
		sort.Strings(sel.Addons)
	}

	if sel.Module == "" && sel.Name != "" {
		sel.Module = fmt.Sprintf("github.com/example/%s", sel.Name)
	}

	return nil
}

// Validate enforces structural and registry-membership rules.
func Validate(sel *Selection, reg *registry.Registry) error {
	if err := ValidateName(sel.Name); err != nil {
		return fmt.Errorf("invalid name %q: %w", sel.Name, err)
	}

	if sel.Module == "" {
		return fmt.Errorf("--module is required (or omit and supply --name)")
	}
	if err := module.CheckPath(sel.Module); err != nil {
		return fmt.Errorf("invalid module path %q: %w", sel.Module, err)
	}

	switch sel.Git {
	case "init", "commit", "none", "":
	default:
		return fmt.Errorf("invalid --git=%s (want init, commit, none)", sel.Git)
	}

	if reg != nil {
		if !reg.Has(registry.KindFramework, sel.Framework) {
			return fmt.Errorf("unknown framework %q", sel.Framework)
		}
		for _, id := range sel.DBs {
			if !reg.Has(registry.KindDB, id) {
				return fmt.Errorf("unknown db %q", id)
			}
		}
		for _, id := range sel.SDKs {
			if !reg.Has(registry.KindSDK, id) {
				return fmt.Errorf("unknown sdk %q", id)
			}
		}
		for _, id := range sel.Addons {
			if !reg.Has(registry.KindAddon, id) {
				return fmt.Errorf("unknown addon %q", id)
			}
		}
		for _, id := range sel.Patterns {
			if !reg.Has(registry.KindPattern, id) {
				return fmt.Errorf("unknown pattern %q", id)
			}
		}
	}
	return nil
}

func dedupSorted(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
