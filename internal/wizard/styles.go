package wizard

import "github.com/charmbracelet/lipgloss"

// ─── Color palette ──────────────────────────────────────────────────────────

var (
	colorPrimary = lipgloss.Color("#01FAC6")
	colorAccent  = lipgloss.Color("205")
	colorDesc    = lipgloss.Color("#40BDA3")
	colorDim     = lipgloss.Color("240")
	colorHint    = lipgloss.Color("190")
	colorSuccess = lipgloss.Color("82")
	colorError   = lipgloss.Color("#FF8700")
	colorSubtle  = lipgloss.Color("238")
)

// ─── Styles ─────────────────────────────────────────────────────────────────

var (
	styleTitle    = lipgloss.NewStyle().Bold(true).Foreground(colorPrimary)
	styleSelected = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	styleDim      = lipgloss.NewStyle().Foreground(colorDim)
	styleHint     = lipgloss.NewStyle().Foreground(colorHint).Italic(true)
	styleCheck    = lipgloss.NewStyle().Foreground(colorSuccess)
	styleCursor   = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	styleError    = lipgloss.NewStyle().Foreground(colorError).Bold(true)
	styleDesc     = lipgloss.NewStyle().Foreground(colorDesc)
	styleBorder   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSubtle).
			Padding(1, 2)
	styleProgressDone = lipgloss.NewStyle().Foreground(colorPrimary)
	styleProgressTodo = lipgloss.NewStyle().Foreground(colorSubtle)
	styleSummaryLabel = lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Width(14)
	styleSummaryValue = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

// ─── Logo ───────────────────────────────────────────────────────────────────

const logoRaw = `
   ██████╗  ██████╗ ██████╗ ██████╗
  ██╔════╝ ██╔═══██╗██╔══██╗██╔══██╗
  ██║  ███╗██║   ██║██████╔╝██████╔╝
  ██║   ██║██║   ██║██╔══██╗██╔═══╝
  ╚██████╔╝╚██████╔╝██████╔╝██║
   ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝`

func renderLogo() string {
	return lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render(logoRaw) + "\n"
}
