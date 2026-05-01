package wizard

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/mod/module"

	"github.com/tienanhnguyen999/gobp/internal/selection"
)

// ─── Stage constants ──────────────────────────────────────────────────────────

type stage int

const (
	stageName stage = iota
	stageModule
	stageFramework
	stageDatabases
	stageSDKs
	stagePatterns
	stageAddons
	stageWebSocket
	stageConfirm
	stageDone
)

const totalStages = 9

// ─── Model ───────────────────────────────────────────────────────────────────

// Model is the top-level Bubble Tea model for the wizard.
type Model struct {
	stage    stage
	aborted  bool
	inputErr string // validation error shown below the active text input

	// text input stages
	nameInput   textinput.Model
	moduleInput textinput.Model

	// list stages: cursor + selected set
	frameworkCursor int
	dbCursor        int
	dbSelected      map[string]bool
	sdkCursor       int
	sdkSelected     map[string]bool
	patternCursor   int
	patternSelected map[string]bool
	addonCursor     int
	addonSelected   map[string]bool

	// websocket toggle
	websocket bool
}

// New creates a Model pre-filled from any flags already provided.
func New(pre selection.Selection) *Model {
	nameIn := textinput.New()
	nameIn.Placeholder = "my-service"
	nameIn.Focus()
	if pre.Name != "" {
		nameIn.SetValue(pre.Name)
	}

	modIn := textinput.New()
	modIn.Placeholder = "github.com/" + gitUser() + "/" + pre.Name
	if pre.Module != "" {
		modIn.SetValue(pre.Module)
	}

	frameworkIdx := 0
	for i, f := range frameworkOptions {
		if f.ID == pre.Framework {
			frameworkIdx = i
			break
		}
	}

	dbSel := map[string]bool{}
	for _, d := range pre.DBs {
		dbSel[d] = true
	}
	sdkSel := map[string]bool{}
	for _, s := range pre.SDKs {
		sdkSel[s] = true
	}
	patSel := map[string]bool{}
	for _, p := range pre.Patterns {
		patSel[p] = true
	}
	addSel := map[string]bool{}
	for _, a := range pre.Addons {
		addSel[a] = true
	}

	// Advance past pre-filled stages
	startStage := stageName
	if pre.Name != "" {
		startStage = stageModule
	}

	m := &Model{
		stage:           startStage,
		nameInput:       nameIn,
		moduleInput:     modIn,
		frameworkCursor: frameworkIdx,
		dbSelected:      dbSel,
		sdkSelected:     sdkSel,
		patternSelected: patSel,
		addonSelected:   addSel,
		websocket:       pre.Websocket,
	}

	// If name was pre-filled, unfocus the name input and focus module input.
	if pre.Name != "" {
		m.nameInput.Blur()
		m.moduleInput.Focus()
	}
	return m
}

// ─── Bubble Tea interface ─────────────────────────────────────────────────────

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate text-input updates.
	var cmd tea.Cmd
	switch m.stage {
	case stageName:
		m.nameInput, cmd = m.nameInput.Update(msg)
	case stageModule:
		m.moduleInput, cmd = m.moduleInput.Update(msg)
	}
	return m, cmd
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.stage {
	case stageName:
		return m.handleText(msg, &m.nameInput)
	case stageModule:
		return m.handleText(msg, &m.moduleInput)
	case stageFramework:
		return m.handleSingleSelect(msg, frameworkOptions, &m.frameworkCursor)
	case stageDatabases:
		return m.handleMultiSelect(msg, dbOptions, &m.dbCursor, m.dbSelected)
	case stageSDKs:
		return m.handleMultiSelect(msg, sdkOptions, &m.sdkCursor, m.sdkSelected)
	case stagePatterns:
		return m.handleMultiSelect(msg, patternOptions, &m.patternCursor, m.patternSelected)
	case stageAddons:
		return m.handleMultiSelect(msg, addonOptions, &m.addonCursor, m.addonSelected)

	case stageWebSocket:
		return m.handleWebSocket(msg)
	case stageConfirm:
		return m.handleConfirm(msg)
	}
	return m, nil
}

// handleText handles key events on a text input stage.
func (m *Model) handleText(msg tea.KeyMsg, input *textinput.Model) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.aborted = true
		return m, tea.Quit
	case tea.KeyEsc:
		if m.stage > stageName {
			m.inputErr = ""
			m.prevStage()
		}
		return m, nil
	case tea.KeyEnter:
		if strings.TrimSpace(input.Value()) == "" && input.Placeholder != "" {
			input.SetValue(input.Placeholder)
		}
		val := strings.TrimSpace(input.Value())
		if val == "" {
			return m, nil
		}
		if err := m.validateStage(val); err != nil {
			m.inputErr = err.Error()
			return m, nil
		}
		m.inputErr = ""
		m.nextStage()
		return m, nil
	}
	m.inputErr = ""
	var cmd tea.Cmd
	*input, cmd = input.Update(msg)
	return m, cmd
}

func (m *Model) handleSingleSelect(msg tea.KeyMsg, opts []Option, cursor *int) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.aborted = true
		return m, tea.Quit
	case tea.KeyLeft, tea.KeyEsc:
		m.prevStage()
	case tea.KeyUp:
		if *cursor > 0 {
			*cursor--
		}
	case tea.KeyDown:
		if *cursor < len(opts)-1 {
			*cursor++
		}
	case tea.KeyEnter:
		m.nextStage()
	case tea.KeyRunes:
		switch msg.Runes[0] {
		case 'k':
			if *cursor > 0 {
				*cursor--
			}
		case 'j':
			if *cursor < len(opts)-1 {
				*cursor++
			}
		}
	}
	return m, nil
}

func (m *Model) handleMultiSelect(msg tea.KeyMsg, opts []Option, cursor *int, selected map[string]bool) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.aborted = true
		return m, tea.Quit
	case tea.KeyLeft, tea.KeyEsc:
		m.prevStage()
	case tea.KeyUp:
		if *cursor > 0 {
			*cursor--
		}
	case tea.KeyDown:
		if *cursor < len(opts)-1 {
			*cursor++
		}
	case tea.KeySpace:
		key := opts[*cursor].ID
		selected[key] = !selected[key]
	case tea.KeyEnter:
		m.nextStage()
	case tea.KeyRunes:
		switch msg.Runes[0] {
		case ' ':
			// Some terminals (e.g. Windows conhost) deliver space as KeyRunes
			// instead of KeySpace.
			key := opts[*cursor].ID
			selected[key] = !selected[key]
		case 'k':
			if *cursor > 0 {
				*cursor--
			}
		case 'j':
			if *cursor < len(opts)-1 {
				*cursor++
			}
		case 's':
			// skip — clear selections and advance
			for k := range selected {
				delete(selected, k)
			}
			m.nextStage()
		}
	}
	return m, nil
}

// handleWebSocket handles the websocket yes/no stage.
func (m *Model) handleWebSocket(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.aborted = true
		return m, tea.Quit
	case tea.KeyLeft, tea.KeyEsc:
		m.prevStage()
	case tea.KeyEnter:
		m.nextStage()
	case tea.KeyRunes:
		switch msg.Runes[0] {
		case 'y', 'Y':
			m.websocket = true
			m.nextStage()
		case 'n', 'N':
			m.websocket = false
			m.nextStage()
		case 'k', 'h', 'j', 'l':
			m.websocket = !m.websocket
		}
	case tea.KeyRight, tea.KeyUp, tea.KeyDown:
		m.websocket = !m.websocket
	}
	return m, nil
}

// handleConfirm handles the final confirm stage.
func (m *Model) handleConfirm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.aborted = true
		return m, tea.Quit
	case tea.KeyLeft, tea.KeyEsc:
		m.prevStage()
		return m, nil
	case tea.KeyEnter:
		m.stage = stageDone
		return m, tea.Quit
	case tea.KeyRunes:
		switch msg.Runes[0] {
		case 'q':
			m.aborted = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

// nextStage advances to the next stage, updating module placeholder if name was just set.
func (m *Model) nextStage() {
	if m.stage == stageName {
		name := strings.TrimSpace(m.nameInput.Value())
		if m.moduleInput.Placeholder == "" || !strings.Contains(m.moduleInput.Placeholder, "/") {
			m.moduleInput.Placeholder = "github.com/" + gitUser() + "/" + name
		}
		m.nameInput.Blur()
		m.moduleInput.Focus()
	}
	if m.stage == stageModule {
		if strings.TrimSpace(m.moduleInput.Value()) == "" {
			m.moduleInput.SetValue(m.moduleInput.Placeholder)
		}
		m.moduleInput.Blur()
	}
	m.stage++
}

func (m *Model) prevStage() {
	if m.stage > stageName {
		m.stage--
	}
	if m.stage == stageName {
		m.nameInput.Focus()
	}
	if m.stage == stageModule {
		m.moduleInput.Focus()
	}
}

func (m *Model) ToSelection() selection.Selection {
	module := strings.TrimSpace(m.moduleInput.Value())
	if module == "" {
		module = m.moduleInput.Placeholder
	}
	return selection.Selection{
		Name:      strings.TrimSpace(m.nameInput.Value()),
		Module:    module,
		Framework: frameworkOptions[m.frameworkCursor].ID,
		DBs:       selectedKeys(m.dbSelected, dbOptions),
		SDKs:      selectedKeys(m.sdkSelected, sdkOptions),
		Patterns:  selectedKeys(m.patternSelected, patternOptions),
		Addons:    selectedKeys(m.addonSelected, addonOptions),
		Websocket: m.websocket,
		Git:       "init",
	}
}

func selectedKeys(m map[string]bool, order []Option) []string {
	var out []string
	for _, o := range order {
		if m[o.ID] {
			out = append(out, o.ID)
		}
	}
	return out
}

// gitUser returns the git config user.name or os username as fallback.
func gitUser() string {
	out, err := exec.Command("git", "config", "user.name").Output()
	if err == nil {
		if name := strings.TrimSpace(string(out)); name != "" {
			// collapse spaces → hyphens, lowercase
			return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
		}
	}
	return "you"
}

// validateStage checks the current text-input value against stage-specific rules.
func (m *Model) validateStage(val string) error {
	switch m.stage {
	case stageName:
		return validateName(val)
	case stageModule:
		return validateModule(val)
	}
	return nil
}

func validateName(val string) error {
	return selection.ValidateName(val)
}

func validateModule(val string) error {
	if err := module.CheckPath(val); err != nil {
		return fmt.Errorf("invalid module path: %w", err)
	}
	return nil
}

// Run launches the TUI and returns the final Selection or an error if aborted.
func Run(pre selection.Selection) (selection.Selection, error) {
	m := New(pre)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return selection.Selection{}, err
	}
	result := final.(*Model)
	if result.aborted {
		return selection.Selection{}, fmt.Errorf("aborted")
	}
	return result.ToSelection(), nil
}
