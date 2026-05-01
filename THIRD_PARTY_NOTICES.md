# Third-party notices

This file credits the projects whose ideas, conventions, or code shaped `gobp`.

## Inspiration

- **[Melkeydev/go-blueprint](https://github.com/Melkeydev/go-blueprint)** (MIT) —
  the overall idea of a Cobra + Bubble Tea Go-project generator with embedded
  templates. `gobp` differs in scope (multi-select features, manifest-driven
  composition, opinionated error catalog, SDK stubs) but go-blueprint was the
  direct prior art. No code is currently copied verbatim; if any snippet is
  borrowed during future work, it will be attributed inline.

## Direct dependencies (gobp itself)

| Project | License |
|---|---|
| github.com/spf13/cobra | Apache-2.0 |
| github.com/charmbracelet/bubbletea | MIT |
| github.com/charmbracelet/bubbles | MIT |
| github.com/charmbracelet/lipgloss | MIT |
| gopkg.in/yaml.v3 | MIT + Apache-2.0 |
| golang.org/x/mod | BSD-3-Clause |
| github.com/google/uuid | BSD-3-Clause |

## Generated-project dependencies

The libraries pinned into the projects `gobp` generates are listed in their
respective `go.mod` files and credited via go.mod / go.sum upstream metadata.
