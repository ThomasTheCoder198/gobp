package render

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"text/template"

	"github.com/tienanhnguyen999/gobp/internal/selection"
	"github.com/tienanhnguyen999/gobp/internal/util"
)

// templatesFS is the embedded tree of *.tmpl files. The actual embed directive
// lives in templates.go (one per templates root) so this file stays generic.

// Ctx is the data shape passed to every template.
type Ctx struct {
	Name      string
	Module    string
	Framework string
	DBs       []string
	SDKs      []string
	Addons    []string
	Patterns  []string
	Websocket bool
	GoVersion string
	Year      int

	// Convenience: present DBs as a map for quick lookup in templates.
	HasDB     map[string]bool
	HasSDK    map[string]bool
	HasAddon  map[string]bool
	HasPattern map[string]bool
}

// NewCtx derives the template context from a Selection.
func NewCtx(sel selection.Selection) Ctx {
	return Ctx{
		Name:       sel.Name,
		Module:     sel.Module,
		Framework:  sel.Framework,
		DBs:        sel.DBs,
		SDKs:       sel.SDKs,
		Addons:     sel.Addons,
		Patterns:   sel.Patterns,
		Websocket:  sel.Websocket,
		GoVersion:  sel.GoVersion,
		Year:       sel.Year,
		HasDB:      toSet(sel.DBs),
		HasSDK:     toSet(sel.SDKs),
		HasAddon:   toSet(sel.Addons),
		HasPattern: toSet(sel.Patterns),
	}
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

// renderTemplate parses the named template from fsys and executes it with ctx.
func renderTemplate(fsys fs.FS, name string, ctx Ctx) ([]byte, error) {
	data, err := fs.ReadFile(fsys, path.Join("templates", name))
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", name, err)
	}
	t, err := template.New(name).Funcs(funcMap()).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse template %s: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return nil, fmt.Errorf("execute template %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"title": strings.Title, //nolint:staticcheck // simple per-word title is fine here
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"quote": func(s string) string { return fmt.Sprintf("%q", s) },
		"hasString": util.StringIn,
		"join": strings.Join,
	}
}

// ensure embed.FS keeps importing in a valid way even if not yet referenced.
var _ = embed.FS{}
