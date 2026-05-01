package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/tienanhnguyen999/gobp/internal/postgen"
	"github.com/tienanhnguyen999/gobp/internal/registry"
	"github.com/tienanhnguyen999/gobp/internal/render"
	"github.com/tienanhnguyen999/gobp/internal/selection"
	"github.com/tienanhnguyen999/gobp/internal/wizard"
)

func newNewCmd() *cobra.Command {
	var (
		name      string
		module    string
		framework string
		dbs       []string
		sdks      []string
		addons    []string
		patterns  []string
		websocket bool
		gitMode   string
		noTidy    bool
		dryRun    bool
		outDir    string
	)

	cmd := &cobra.Command{
		Use:   "new",
		Short: "Generate a new Go project",
		RunE: func(cmd *cobra.Command, _ []string) error {
			sel := selection.Selection{
				Name:      name,
				Module:    module,
				Framework: framework,
				DBs:       dbs,
				SDKs:      sdks,
				Addons:    addons,
				Patterns:  patterns,
				Websocket: websocket,
				Git:       gitMode,
				NoTidy:    noTidy,
				DryRun:    dryRun,
			}

			// Launch interactive wizard when --name was not provided.
			if !cmd.Flags().Changed("name") {
				var err error
				sel, err = wizard.Run(sel)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
			}

			reg, err := registry.Load()
			if err != nil {
				return fmt.Errorf("load registry: %w", err)
			}

			if err := selection.Resolve(&sel, reg); err != nil {
				return fmt.Errorf("resolve selection: %w", err)
			}
			if err := selection.Validate(&sel, reg); err != nil {
				return fmt.Errorf("validate selection: %w", err)
			}

			plan, err := render.BuildPlan(sel, reg)
			if err != nil {
				return fmt.Errorf("build plan: %w", err)
			}

			if dryRun {
				return render.PrintPlan(os.Stdout, plan)
			}

			target := outDir
			if target == "" {
				target = sel.Name
			}
			if !filepath.IsAbs(target) {
				wd, _ := os.Getwd()
				target = filepath.Join(wd, target)
			}

			if err := RunWithSpinner("Generating your project...", func() error {
				if err := render.Execute(plan, target); err != nil {
					return fmt.Errorf("render: %w", err)
				}
				if err := postgen.Run(target, sel); err != nil {
					return fmt.Errorf("post-gen: %w", err)
				}
				return nil
			}); err != nil {
				return err
			}

			fmt.Fprintf(os.Stdout, "\n  Generated %s at %s\n", sel.Name, target)
			fmt.Fprintf(os.Stdout, "  Next: cd %s && go build ./...\n\n", sel.Name)
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "project name (launches wizard if omitted)")
	cmd.Flags().StringVar(&module, "module", "", "Go module path (defaults to github.com/<git-user>/<name>)")
	cmd.Flags().StringVar(&framework, "framework", "gin", "HTTP framework: gin, echo, fiber")
	cmd.Flags().StringSliceVar(&dbs, "db", nil, "databases (multi): postgres, mysql, sqlite, mongo, redis, cassandra")
	cmd.Flags().StringSliceVar(&sdks, "sdk", nil, "SDK clients (multi): openai, stripe")
	cmd.Flags().StringSliceVar(&addons, "addon", nil, "addons (multi): docker, githubactions")
	cmd.Flags().StringSliceVar(&patterns, "pattern", nil, "patterns (multi): worker")
	cmd.Flags().BoolVar(&websocket, "websocket", false, "scaffold a websocket handler")
	cmd.Flags().StringVar(&gitMode, "git", "init", "git mode: init, commit, none")
	cmd.Flags().BoolVar(&noTidy, "no-tidy", false, "skip go mod tidy after generation")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "print render plan without writing")
	cmd.Flags().StringVar(&outDir, "out", "", "output directory (default: ./<name>)")

	return cmd
}
