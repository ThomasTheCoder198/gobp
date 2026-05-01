package wizard

import (
	"fmt"
	"strings"

	"github.com/tienanhnguyen999/gobp/internal/selection"
)

// ─── Stage labels ───────────────────────────────────────────────────────────

func stageLabel(s stage) string {
	switch s {
	case stageName:
		return "Project name"
	case stageModule:
		return "Go module path"
	case stageFramework:
		return "Framework"
	case stageDatabases:
		return "Databases"
	case stageSDKs:
		return "SDK clients"
	case stagePatterns:
		return "Patterns"
	case stageAddons:
		return "Addons"
	case stageWebSocket:
		return "WebSocket"
	case stageConfirm:
		return "Confirm"
	default:
		return ""
	}
}

// ─── Progress bar ───────────────────────────────────────────────────────────

func renderProgress(current stage) string {
	step := int(current) + 1
	var bar strings.Builder
	for i := 1; i <= totalStages; i++ {
		if i <= step {
			bar.WriteString(styleProgressDone.Render("█"))
		} else {
			bar.WriteString(styleProgressTodo.Render("░"))
		}
	}
	label := stageLabel(current)
	return fmt.Sprintf("  %s  %s  %s",
		styleTitle.Render(fmt.Sprintf("Step %d of %d", step, totalStages)),
		bar.String(),
		styleDesc.Render(label),
	)
}

// ─── Header ─────────────────────────────────────────────────────────────────

func renderHeader(current stage) string {
	var b strings.Builder
	b.WriteString(renderLogo())
	b.WriteString(renderProgress(current))
	b.WriteString("\n")
	b.WriteString(styleDim.Render("  " + strings.Repeat("─", 46)))
	b.WriteString("\n\n")
	return b.String()
}

// ─── Hint bar ───────────────────────────────────────────────────────────────

func renderHint(parts ...string) string {
	return styleHint.Render("  " + strings.Join(parts, " · "))
}

// ─── View ───────────────────────────────────────────────────────────────────

func (m *Model) View() string {
	if m.stage == stageDone {
		return ""
	}

	var b strings.Builder
	b.WriteString(renderHeader(m.stage))

	switch m.stage {
	case stageName:
		b.WriteString(styleSelected.Render("  Project name") + "\n\n")
		b.WriteString("  " + m.nameInput.View() + "\n")
		if m.inputErr != "" {
			b.WriteString(styleError.Render("  ✗ "+m.inputErr) + "\n")
		}
		b.WriteString("\n" + renderHint("enter to confirm") + "\n")

	case stageModule:
		b.WriteString(styleSelected.Render("  Go module path") + "\n\n")
		b.WriteString("  " + m.moduleInput.View() + "\n")
		if m.inputErr != "" {
			b.WriteString(styleError.Render("  ✗ "+m.inputErr) + "\n")
		}
		b.WriteString("\n" + renderHint("enter to confirm", "← back") + "\n")

	case stageFramework:
		b.WriteString(styleSelected.Render("  Framework") + "\n\n")
		for i, f := range frameworkOptions {
			if i == m.frameworkCursor {
				b.WriteString("  " + styleCursor.Render("▸ ") + styleSelected.Render(f.ID))
				b.WriteString(styleDesc.Render(" — "+f.Desc) + "\n")
			} else {
				b.WriteString("    " + styleDim.Render(f.ID))
				b.WriteString(styleDim.Render(" — "+f.Desc) + "\n")
			}
		}
		b.WriteString("\n" + renderHint("↑↓/jk move", "enter confirm", "← back") + "\n")

	case stageDatabases:
		b.WriteString(styleSelected.Render("  Databases") + styleDim.Render("  multi-select") + "\n\n")
		b.WriteString(renderOptionList(dbOptions, m.dbCursor, m.dbSelected))
		b.WriteString("\n" + renderHint("↑↓/jk move", "space toggle", "enter confirm", "s skip", "← back") + "\n")

	case stageSDKs:
		b.WriteString(styleSelected.Render("  SDK clients") + styleDim.Render("  optional") + "\n\n")
		b.WriteString(renderOptionList(sdkOptions, m.sdkCursor, m.sdkSelected))
		b.WriteString("\n" + renderHint("↑↓/jk move", "space toggle", "enter confirm", "s skip", "← back") + "\n")

	case stagePatterns:
		b.WriteString(styleSelected.Render("  Patterns") + styleDim.Render("  optional") + "\n\n")
		b.WriteString(renderOptionList(patternOptions, m.patternCursor, m.patternSelected))
		b.WriteString("\n" + renderHint("↑↓/jk move", "space toggle", "enter confirm", "s skip", "← back") + "\n")

	case stageAddons:
		b.WriteString(styleSelected.Render("  Addons") + styleDim.Render("  optional") + "\n\n")
		b.WriteString(renderOptionList(addonOptions, m.addonCursor, m.addonSelected))
		b.WriteString("\n" + renderHint("↑↓/jk move", "space toggle", "enter confirm", "s skip", "← back") + "\n")

	case stageWebSocket:
		b.WriteString(styleSelected.Render("  WebSocket handler?") + "\n\n")
		yes, no := styleDim, styleDim
		if m.websocket {
			yes = styleSelected
		} else {
			no = styleSelected
		}
		b.WriteString("    " + yes.Render("[ Yes ]") + "   " + no.Render("[ No ]") + "\n")
		b.WriteString("\n" + renderHint("y/n toggle", "enter confirm", "← back") + "\n")

	case stageConfirm:
		sel := m.ToSelection()
		b.WriteString(styleSelected.Render("  Ready to generate") + "\n\n")
		b.WriteString(styleBorder.Render(renderSummary(sel)) + "\n\n")
		b.WriteString(styleDim.Render("  Equivalent command:") + "\n")
		b.WriteString(styleHint.Render("  "+flagCommand(sel)) + "\n")
		b.WriteString("\n" + renderHint("enter generate", "← back", "q abort") + "\n")
	}

	return b.String()
}

// ─── Render helpers ─────────────────────────────────────────────────────────

func renderOptionList(opts []Option, cursor int, selected map[string]bool) string {
	var b strings.Builder
	for i, o := range opts {
		var mark string
		if selected[o.ID] {
			mark = styleCheck.Render("[✓]")
		} else {
			mark = styleDim.Render("[ ]")
		}
		if i == cursor {
			b.WriteString("  " + styleCursor.Render("▸ ") + mark + " " + styleSelected.Render(o.ID))
			b.WriteString(styleDesc.Render(" — "+o.Desc) + "\n")
		} else {
			b.WriteString("    " + mark + " " + styleDim.Render(o.ID))
			b.WriteString(styleDim.Render(" — "+o.Desc) + "\n")
		}
	}
	return b.String()
}

func renderSummary(sel selection.Selection) string {
	var b strings.Builder
	line := func(label, val string) {
		b.WriteString(styleSummaryLabel.Render(label) + styleSummaryValue.Render(val) + "\n")
	}
	line("Name", sel.Name)
	line("Module", sel.Module)
	line("Framework", sel.Framework)
	if len(sel.DBs) > 0 {
		line("Databases", strings.Join(sel.DBs, ", "))
	}
	if len(sel.SDKs) > 0 {
		line("SDKs", strings.Join(sel.SDKs, ", "))
	}
	if len(sel.Patterns) > 0 {
		line("Patterns", strings.Join(sel.Patterns, ", "))
	}
	if len(sel.Addons) > 0 {
		line("Addons", strings.Join(sel.Addons, ", "))
	}
	if sel.Websocket {
		line("WebSocket", "yes")
	}
	return strings.TrimRight(b.String(), "\n")
}

func flagCommand(sel selection.Selection) string {
	parts := []string{"gobp new", "--name " + sel.Name, "--module " + sel.Module}
	if sel.Framework != "gin" {
		parts = append(parts, "--framework "+sel.Framework)
	}
	if len(sel.DBs) > 0 {
		parts = append(parts, "--db "+strings.Join(sel.DBs, ","))
	}
	if len(sel.SDKs) > 0 {
		parts = append(parts, "--sdk "+strings.Join(sel.SDKs, ","))
	}
	if len(sel.Patterns) > 0 {
		parts = append(parts, "--pattern "+strings.Join(sel.Patterns, ","))
	}
	if len(sel.Addons) > 0 {
		parts = append(parts, "--addon "+strings.Join(sel.Addons, ","))
	}
	if sel.Websocket {
		parts = append(parts, "--websocket")
	}
	return strings.Join(parts, " ")
}
