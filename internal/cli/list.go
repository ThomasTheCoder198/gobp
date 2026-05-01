package cli

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/tienanhnguyen999/gobp/internal/registry"
)

// kindOrder controls the display order of sections in `gobp list`.
var kindOrder = []registry.Kind{
	registry.KindFramework,
	registry.KindDB,
	registry.KindSDK,
	registry.KindAddon,
	registry.KindPattern,
	registry.KindLogger,
}

var kindLabel = map[registry.Kind]string{
	registry.KindFramework: "Frameworks",
	registry.KindDB:        "Databases",
	registry.KindSDK:       "SDKs",
	registry.KindAddon:     "Addons",
	registry.KindPattern:   "Patterns",
	registry.KindLogger:    "Loggers",
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [kind]",
		Short: "List available features (kind: dbs, frameworks, sdks, addons, patterns, features)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			reg, err := registry.Load()
			if err != nil {
				return err
			}
			kind := "features"
			if len(args) > 0 {
				kind = args[0]
			}

			switch kind {
			case "dbs", "db":
				printSection(reg.ByKind(registry.KindDB), "")
			case "frameworks", "framework":
				printSection(reg.ByKind(registry.KindFramework), "")
			case "sdks", "sdk":
				printSection(reg.ByKind(registry.KindSDK), "")
			case "addons", "addon":
				printSection(reg.ByKind(registry.KindAddon), "")
			case "patterns", "pattern":
				printSection(reg.ByKind(registry.KindPattern), "")
			case "features":
				first := true
				for _, k := range kindOrder {
					ms := reg.ByKind(k)
					if len(ms) == 0 {
						continue
					}
					if !first {
						fmt.Println()
					}
					first = false
					label := kindLabel[k]
					fmt.Printf("%s\n", label)
					printSection(ms, "  ")
				}
			default:
				return fmt.Errorf("unknown kind: %s (want: dbs, frameworks, sdks, addons, patterns, features)", kind)
			}
			return nil
		},
	}
	return cmd
}

func printSection(ms []registry.Manifest, indent string) {
	sort.Slice(ms, func(i, j int) bool { return ms[i].ID < ms[j].ID })
	for _, m := range ms {
		fmt.Printf("%s%-15s %s\n", indent, m.ID, m.Display)
	}
}
