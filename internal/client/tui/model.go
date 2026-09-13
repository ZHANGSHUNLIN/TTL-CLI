package tui

import (
	"fmt"
	"io"
	"strings"

	clientapp "ttl-cli/internal/client/app"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const narrowWidth = 72

type ResourceService interface {
	ListResources() ([]clientapp.Resource, error)
	FindResources(string) ([]clientapp.Resource, error)
	CreateResource(string, string, []string) (clientapp.Resource, error)
	UpdateResourceValue(string, string) (clientapp.Resource, error)
	AddResourceTags(string, []string) (clientapp.Resource, error)
	DeleteResourceTag(string, string) (clientapp.Resource, error)
	DeleteResourceWithCleanup(string) (clientapp.DeleteResult, error)
}

type RunOptions struct {
	In            io.Reader
	Out           io.Writer
	InitialWidth  int
	InitialHeight int
}

type screen int

const (
	browseScreen screen = iota
	detailScreen
	searchScreen
	createScreen
	editScreen
	tagScreen
	deleteScreen
	discardScreen
	helpScreen
)

type loadMsg struct {
	query     string
	resources []clientapp.Resource
	err       error
}

type mutationMsg struct {
	warning string
	err     error
}

type Model struct {
	service      ResourceService
	resources    []clientapp.Resource
	selected     int
	query        string
	screen       screen
	previous     screen
	width        int
	height       int
	loading      bool
	busy         bool
	status       string
	err          error
	search       textinput.Model
	key          textinput.Model
	value        textarea.Model
	tags         textinput.Model
	tag          textinput.Model
	dirty        bool
	detailOffset int
}

func NewModel(service ResourceService, width, height int) Model {
	search := textinput.New()
	search.Prompt = "/ "
	search.Placeholder = "search key, value, or tag"
	key := textinput.New()
	key.Prompt = "Key: "
	value := textarea.New()
	value.Prompt = "│ "
	value.Placeholder = "Resource content"
	value.SetWidth(60)
	value.SetHeight(8)
	tags := textinput.New()
	tags.Prompt = "Tags: "
	tags.Placeholder = "comma,separated"
	tag := textinput.New()
	tag.Prompt = "Tags (+name or -name): "
	model := Model{
		service: service, width: width, height: height, screen: browseScreen, loading: true,
		search: search, key: key, value: value, tags: tags, tag: tag,
	}
	model.resizeInputs()
	return model
}

func Run(service ResourceService, opts RunOptions) (err error) {
	if service == nil {
		return fmt.Errorf("TUI service is not initialized")
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("TUI crashed: %v", recovered)
		}
	}()
	programOptions := []tea.ProgramOption{tea.WithAltScreen()}
	if opts.In != nil {
		programOptions = append(programOptions, tea.WithInput(opts.In))
	}
	if opts.Out != nil {
		programOptions = append(programOptions, tea.WithOutput(opts.Out))
	}
	_, err = tea.NewProgram(NewModel(service, opts.InitialWidth, opts.InitialHeight), programOptions...).Run()
	return err
}

func (m Model) Init() tea.Cmd { return m.loadResourcesCmd("") }

func (m Model) loadResourcesCmd(query string) tea.Cmd {
	return func() tea.Msg {
		var resources []clientapp.Resource
		var err error
		if query == "" {
			resources, err = m.service.ListResources()
		} else {
			resources, err = m.service.FindResources(query)
			if kind, ok := clientapp.ErrorKindOf(err); ok && kind == clientapp.ErrorNotFound {
				err = nil
			}
		}
		return loadMsg{query: query, resources: resources, err: err}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.resizeInputs()
		return m, nil
	case loadMsg:
		if msg.query != m.query {
			return m, nil
		}
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.status = "Unable to read resources; press r to retry"
			return m, nil
		}
		m.err = nil
		m.resources = msg.resources
		if len(m.resources) == 0 {
			m.selected = 0
		} else if m.selected >= len(m.resources) {
			m.selected = len(m.resources) - 1
		}
		return m, nil
	case mutationMsg:
		m.busy = false
		if msg.err != nil {
			m.err = msg.err
			m.status = "Operation failed; your input was kept"
			return m, nil
		}
		m.err = nil
		m.dirty = false
		m.status = "Saved"
		if msg.warning != "" {
			m.status = "Saved with warning: " + msg.warning
		}
		m.screen = browseScreen
		return m, m.loadResourcesCmd(m.query)
	}

	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch m.screen {
	case searchScreen:
		return m.updateSearch(key)
	case createScreen, editScreen:
		return m.updateEditor(key)
	case tagScreen:
		return m.updateTag(key)
	case deleteScreen:
		return m.updateDelete(key)
	case discardScreen:
		return m.updateDiscard(key)
	case detailScreen:
		switch key.String() {
		case "up", "k":
			if m.detailOffset > 0 {
				m.detailOffset--
			}
			return m, nil
		case "down", "j":
			m.detailOffset++
			return m, nil
		case "esc", "q":
			m.screen = browseScreen
			return m, nil
		case "e":
			m.startEdit()
			return m, nil
		case "t":
			m.screen = tagScreen
			m.tag.SetValue("")
			m.tag.Focus()
			return m, nil
		case "d":
			m.screen = deleteScreen
			return m, nil
		}
	case helpScreen:
		if key.String() == "esc" || key.String() == "q" || m.screen == helpScreen && key.String() == "?" {
			m.screen = browseScreen
			return m, nil
		}
	}
	return m.updateBrowse(key)
}

func (m Model) updateBrowse(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
	case "down", "j":
		if m.selected+1 < len(m.resources) {
			m.selected++
		}
	case "enter":
		if len(m.resources) > 0 {
			m.screen = detailScreen
			m.detailOffset = 0
		}
	case "/":
		m.screen = searchScreen
		m.search.SetValue(m.query)
		m.search.Focus()
	case "n":
		m.startCreate()
	case "e":
		m.startEdit()
	case "t":
		if len(m.resources) > 0 {
			m.screen = tagScreen
			m.tag.SetValue("")
			m.tag.Focus()
		}
	case "d":
		if len(m.resources) > 0 {
			m.screen = deleteScreen
		}
	case "?":
		m.screen = helpScreen
	case "r":
		m.loading = true
		return m, m.loadResourcesCmd(m.query)
	case "esc":
		if m.query != "" {
			m.query = ""
			m.search.SetValue("")
			m.loading = true
			return m, m.loadResourcesCmd("")
		}
	}
	return m, nil
}

func (m Model) updateSearch(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "esc":
		m.query = ""
		m.search.SetValue("")
		m.screen = browseScreen
		return m, m.loadResourcesCmd("")
	case "enter":
		m.query = strings.TrimSpace(m.search.Value())
		m.screen = browseScreen
		return m, m.loadResourcesCmd(m.query)
	}
	var cmd tea.Cmd
	m.search, cmd = m.search.Update(key)
	m.query = strings.TrimSpace(m.search.Value())
	return m, tea.Batch(cmd, m.loadResourcesCmd(m.query))
}

func (m *Model) startCreate() {
	m.screen = createScreen
	m.key.SetValue("")
	m.value.SetValue("")
	m.tags.SetValue("")
	m.key.Focus()
	m.value.Blur()
	m.tags.Blur()
	m.dirty = false
}

func (m *Model) startEdit() {
	if len(m.resources) == 0 {
		return
	}
	m.screen = editScreen
	m.value.SetValue(m.resources[m.selected].Value.Val)
	m.value.Focus()
	m.dirty = false
}

func (m Model) updateEditor(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.String() == "ctrl+s" {
		if m.busy {
			return m, nil
		}
		if m.screen == createScreen && strings.TrimSpace(m.key.Value()) == "" {
			m.status = "Key is required"
			return m, nil
		}
		m.busy = true
		if m.screen == createScreen {
			keyValue, value, tags := strings.TrimSpace(m.key.Value()), m.value.Value(), splitTags(m.tags.Value())
			return m, func() tea.Msg {
				_, err := m.service.CreateResource(keyValue, value, tags)
				return mutationMsg{err: err}
			}
		}
		resourceKey, value := m.resources[m.selected].Key.Key, m.value.Value()
		return m, func() tea.Msg {
			_, err := m.service.UpdateResourceValue(resourceKey, value)
			return mutationMsg{err: err}
		}
	}
	if key.String() == "esc" || key.String() == "ctrl+c" {
		if m.dirty {
			m.previous = m.screen
			m.screen = discardScreen
			return m, nil
		}
		m.screen = browseScreen
		return m, nil
	}
	if m.screen == createScreen && key.String() == "tab" {
		switch {
		case m.key.Focused():
			m.key.Blur()
			m.value.Focus()
		case m.value.Focused():
			m.value.Blur()
			m.tags.Focus()
		default:
			m.tags.Blur()
			m.key.Focus()
		}
		return m, nil
	}
	var cmd tea.Cmd
	if m.screen == editScreen || m.value.Focused() {
		m.value, cmd = m.value.Update(key)
	} else if m.key.Focused() {
		m.key, cmd = m.key.Update(key)
	} else {
		m.tags, cmd = m.tags.Update(key)
	}
	m.dirty = m.editorDirty()
	return m, cmd
}

func (m Model) editorDirty() bool {
	if m.screen == createScreen {
		return m.key.Value() != "" || m.value.Value() != "" || m.tags.Value() != ""
	}
	return len(m.resources) > 0 && m.value.Value() != m.resources[m.selected].Value.Val
}

func (m Model) updateDiscard(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "y":
		m.dirty = false
		m.screen = browseScreen
		m.status = "Discarded unsaved changes"
	case "n", "esc":
		m.screen = m.previous
	}
	return m, nil
}

func (m Model) updateTag(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.String() == "esc" {
		m.screen = browseScreen
		return m, nil
	}
	if key.String() == "enter" && !m.busy {
		input := strings.TrimSpace(m.tag.Value())
		if input == "" {
			return m, nil
		}
		m.busy = true
		resourceKey := m.resources[m.selected].Key.Key
		return m, func() tea.Msg {
			var err error
			if strings.HasPrefix(input, "-") {
				_, err = m.service.DeleteResourceTag(resourceKey, strings.TrimSpace(strings.TrimPrefix(input, "-")))
			} else {
				_, err = m.service.AddResourceTags(resourceKey, splitTags(strings.TrimPrefix(input, "+")))
			}
			return mutationMsg{err: err}
		}
	}
	var cmd tea.Cmd
	m.tag, cmd = m.tag.Update(key)
	return m, cmd
}

func (m Model) updateDelete(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch key.String() {
	case "n", "esc":
		m.screen = browseScreen
	case "y":
		if m.busy {
			return m, nil
		}
		m.busy = true
		resourceKey := m.resources[m.selected].Key.Key
		return m, func() tea.Msg {
			result, err := m.service.DeleteResourceWithCleanup(resourceKey)
			warnings := []string{}
			if result.HistoryCleanupError != nil {
				warnings = append(warnings, result.HistoryCleanupError.Error())
			}
			if result.AuditCleanupError != nil {
				warnings = append(warnings, result.AuditCleanupError.Error())
			}
			return mutationMsg{warning: strings.Join(warnings, "; "), err: err}
		}
	}
	return m, nil
}

func (m *Model) resizeInputs() {
	width := m.width - 8
	if width < 20 {
		width = 20
	}
	if width > 80 {
		width = 80
	}
	m.key.Width = width
	m.tags.Width = width
	m.tag.Width = width
	m.search.Width = width
	m.value.SetWidth(width)
	height := m.height - 12
	if height < 4 {
		height = 4
	}
	if height > 16 {
		height = 16
	}
	m.value.SetHeight(height)
}

func splitTags(value string) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, raw := range strings.Split(value, ",") {
		tag := strings.TrimSpace(raw)
		if tag != "" && !seen[tag] {
			seen[tag] = true
			result = append(result, tag)
		}
	}
	return result
}

func (m Model) View() string {
	if m.width > 0 && (m.width < 30 || m.height < 8) {
		return "Terminal too small; resize or press q.\n"
	}
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63")).Render("TTL Resources")
	header := fmt.Sprintf("%s  local  query=%q  results=%d", title, m.query, len(m.resources))
	if m.loading {
		return header + "\n\nLoading resources…\n"
	}
	if m.err != nil && len(m.resources) == 0 {
		return header + "\n\nError: " + m.err.Error() + "\n" + m.status + "\n"
	}
	var body string
	switch m.screen {
	case searchScreen:
		body = m.search.View() + "\n\nEnter keep filter • Esc clear"
	case createScreen:
		body = "Create resource\n\n" + m.key.View() + "\nValue:\n" + m.value.View() + "\n" + m.tags.View() + "\n\nCtrl+S save • Tab next • Esc cancel"
	case editScreen:
		body = "Edit " + m.currentKey() + "\n\n" + m.value.View() + "\n\nCtrl+S save • Esc cancel"
	case tagScreen:
		body = "Manage tags for " + m.currentKey() + "\nCurrent: " + strings.Join(m.currentTags(), ", ") + "\n\n" + m.tag.View() + "\nUse name to add or -name to remove • Enter apply • Esc cancel"
	case deleteScreen:
		body = "Delete " + m.currentKey() + "?\n" + truncate(m.currentValue(), 120) + "\n\ny confirm • n/Esc cancel"
	case discardScreen:
		body = "Discard unsaved changes? They cannot be recovered.\n\ny discard • n/Esc continue editing"
	case helpScreen:
		body = "Keys\n  ↑/k ↓/j select   Enter details   / search   n new\n  e edit   t tags   d delete   r retry   q quit"
	case detailScreen:
		body = m.scrollableDetailView() + "\n\n↑/↓ scroll • Esc/q back • e edit • t tags • d delete"
	default:
		body = m.browserView()
	}
	footer := m.status
	if m.err != nil {
		footer = "Error: " + m.err.Error() + " • " + footer
	}
	if m.busy {
		footer = "Working…"
	}
	if footer != "" {
		return clipView(header+"\n\n"+body+"\n\n"+footer+"\n", m.width, m.height)
	}
	return clipView(header+"\n\n"+body+"\n", m.width, m.height)
}

func (m Model) browserView() string {
	if len(m.resources) == 0 {
		if m.query != "" {
			return fmt.Sprintf("No results for %q. Press / to change or Esc to clear.", m.query)
		}
		return "No resources yet. Press n to add the first one."
	}
	lines := make([]string, 0, len(m.resources))
	for i, resource := range m.resources {
		prefix := "  "
		if i == m.selected {
			prefix = "› "
		}
		lines = append(lines, prefix+resource.Key.Key+tagSuffix(resource.Value.Tag))
	}
	list := strings.Join(lines, "\n")
	if m.width > 0 && m.width < narrowWidth {
		return list + "\n\nEnter details • / search • n new • ? help • q quit"
	}
	detailWidth := m.width - 38
	if detailWidth < 20 {
		detailWidth = 20
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, lipgloss.NewStyle().Width(32).Render(list), "  ", lipgloss.NewStyle().Width(detailWidth).Render(m.detailView())) + "\n\n/ search • n new • e edit • t tags • d delete • ? help • q quit"
}

func (m Model) detailView() string {
	if len(m.resources) == 0 {
		return ""
	}
	resource := m.resources[m.selected]
	return fmt.Sprintf("%s\n\n%s\n\nTags: %s\nCreated: %d  Updated: %d", resource.Key.Key, resource.Value.Val, strings.Join(resource.Value.Tag, ", "), resource.Value.CreatedAt, resource.Value.UpdatedAt)
}

func (m Model) scrollableDetailView() string {
	lines := wrapLines(m.detailView(), maxInt(20, m.width-4))
	visible := maxInt(3, m.height-7)
	maxOffset := len(lines) - visible
	if maxOffset < 0 {
		maxOffset = 0
	}
	offset := m.detailOffset
	if offset > maxOffset {
		offset = maxOffset
	}
	end := offset + visible
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[offset:end], "\n")
}

func wrapLines(value string, width int) []string {
	var result []string
	for _, line := range strings.Split(value, "\n") {
		runes := []rune(line)
		if len(runes) == 0 {
			result = append(result, "")
			continue
		}
		for len(runes) > width {
			result = append(result, string(runes[:width]))
			runes = runes[width:]
		}
		result = append(result, string(runes))
	}
	return result
}

func maxInt(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func (m Model) currentKey() string {
	if len(m.resources) == 0 {
		return ""
	}
	return m.resources[m.selected].Key.Key
}

func (m Model) currentValue() string {
	if len(m.resources) == 0 {
		return ""
	}
	return m.resources[m.selected].Value.Val
}

func (m Model) currentTags() []string {
	if len(m.resources) == 0 {
		return nil
	}
	return m.resources[m.selected].Value.Tag
}

func tagSuffix(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	return "  [" + strings.Join(tags, ",") + "]"
}

func truncate(value string, limit int) string {
	value = strings.ReplaceAll(value, "\n", " ")
	if len(value) <= limit {
		return value
	}
	return value[:limit] + "…"
}

func clipView(value string, width, height int) string {
	if width <= 0 || height <= 0 {
		return value
	}
	lines := strings.Split(value, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	for index, line := range lines {
		runes := []rune(line)
		if len(runes) > width {
			lines[index] = string(runes[:width])
		}
	}
	return strings.Join(lines, "\n")
}
