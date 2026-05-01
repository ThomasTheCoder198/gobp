package registry

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed manifests/*.yaml
var manifestsFS embed.FS

// Registry is the indexed set of feature manifests embedded in the binary.
type Registry struct {
	all   []Manifest
	byID  map[string]*Manifest
	byKnd map[Kind][]*Manifest
}

// Load parses every embedded manifest and builds the registry.
func Load() (*Registry, error) {
	return loadFS(manifestsFS, "manifests")
}

func loadFS(fsys fs.FS, dir string) (*Registry, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("read manifests dir: %w", err)
	}
	r := &Registry{
		byID:  map[string]*Manifest{},
		byKnd: map[Kind][]*Manifest{},
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}
		var m Manifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if err := validateManifest(&m); err != nil {
			return nil, fmt.Errorf("invalid manifest %s: %w", e.Name(), err)
		}
		if _, dup := r.byID[m.ID]; dup {
			return nil, fmt.Errorf("duplicate manifest id %q (in %s)", m.ID, e.Name())
		}
		r.all = append(r.all, m)
	}
	// Re-index pointers after the slice is final.
	for i := range r.all {
		m := &r.all[i]
		r.byID[m.ID] = m
		r.byKnd[m.Kind] = append(r.byKnd[m.Kind], m)
	}
	for k := range r.byKnd {
		sort.Slice(r.byKnd[k], func(i, j int) bool { return r.byKnd[k][i].ID < r.byKnd[k][j].ID })
	}
	return r, nil
}

func validateManifest(m *Manifest) error {
	if m.ID == "" {
		return fmt.Errorf("missing id")
	}
	if m.Kind == "" {
		return fmt.Errorf("missing kind")
	}
	switch m.Kind {
	case KindFramework, KindLogger, KindDB, KindSDK, KindAddon, KindPattern:
	default:
		return fmt.Errorf("unknown kind %q", m.Kind)
	}
	if m.Display == "" {
		m.Display = m.ID
	}
	for i, f := range m.Fragments {
		if f.Target == "" {
			return fmt.Errorf("fragment[%d]: missing target", i)
		}
		if f.File == "" && f.Inline == "" {
			return fmt.Errorf("fragment[%d]: must set file or inline", i)
		}
	}
	return nil
}

// Has reports whether a manifest of the given kind and id is registered.
func (r *Registry) Has(kind Kind, id string) bool {
	m := r.byID[id]
	return m != nil && m.Kind == kind
}

// Get returns the manifest with id, or nil.
func (r *Registry) Get(id string) *Manifest {
	return r.byID[id]
}

// ByKind returns all manifests of the given kind, sorted by ID.
func (r *Registry) ByKind(kind Kind) []Manifest {
	ms := r.byKnd[kind]
	out := make([]Manifest, len(ms))
	for i, m := range ms {
		out[i] = *m
	}
	return out
}

// All returns every registered manifest (copy).
func (r *Registry) All() []Manifest {
	out := make([]Manifest, len(r.all))
	copy(out, r.all)
	return out
}
