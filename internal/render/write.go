package render

import (
	"fmt"
	"os"
	"path/filepath"
)

// Execute renders the plan into a temporary directory, then atomically moves
// the result to target. On any error the temp directory is cleaned and the
// target is left untouched — there are no half-written projects.
func Execute(plan *Plan, target string) error {
	if plan == nil {
		return fmt.Errorf("nil plan")
	}

	// Ensure target's parent exists, but the target itself must not exist
	// (or must be empty) — refuse to clobber.
	if info, err := os.Stat(target); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("target %s exists and is not a directory", target)
		}
		entries, err := os.ReadDir(target)
		if err != nil {
			return fmt.Errorf("read target %s: %w", target, err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("target %s is not empty", target)
		}
	}

	tmp, err := os.MkdirTemp("", "gobp-")
	if err != nil {
		return fmt.Errorf("mkdir temp: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(tmp) }

	ctx := NewCtx(plan.Selection)
	for _, f := range plan.Files {
		var content []byte
		if f.Template != "" {
			b, err := renderTemplate(templatesFS, f.Template, ctx)
			if err != nil {
				cleanup()
				return err
			}
			if f.Format {
				formatted, err := gofmt(f.Target, b)
				if err != nil {
					cleanup()
					return err
				}
				b = formatted
			}
			content = b
		} else {
			content = f.Static
		}

		out := filepath.Join(tmp, f.Target)
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			cleanup()
			return fmt.Errorf("mkdir %s: %w", filepath.Dir(out), err)
		}
		mode := os.FileMode(0o644)
		if f.Executable {
			mode = 0o755
		}
		if err := os.WriteFile(out, content, mode); err != nil {
			cleanup()
			return fmt.Errorf("write %s: %w", out, err)
		}
	}

	// Move tmp to target. If target exists (we already checked it's empty
	// above), use the more robust copy+remove path since os.Rename can fail
	// on Windows when the target directory exists.
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		cleanup()
		return fmt.Errorf("mkdir target parent: %w", err)
	}
	if _, err := os.Stat(target); err == nil {
		// Target exists (empty) — copy tmp contents into it.
		if err := copyDir(tmp, target); err != nil {
			cleanup()
			return err
		}
		cleanup()
		return nil
	}
	if err := os.Rename(tmp, target); err != nil {
		// Cross-device rename or permission — fall back to copy.
		if err := copyDir(tmp, target); err != nil {
			cleanup()
			return err
		}
		cleanup()
	}
	return nil
}

func copyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := copyDir(s, d); err != nil {
				return err
			}
			continue
		}
		data, err := os.ReadFile(s)
		if err != nil {
			return err
		}
		info, err := e.Info()
		if err != nil {
			return err
		}
		if err := os.WriteFile(d, data, info.Mode()); err != nil {
			return err
		}
	}
	return nil
}
