package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	// Obtener URL del API desde argumentos
	apiURL := ""
	if len(os.Args) > 1 {
		apiURL = os.Args[1]
	}

	if apiURL == "" {
		fmt.Println("Usage: snag <api-url>")
		fmt.Println("Example: snag http://localhost:8000")
		os.Exit(1)
	}

	// Create initial model
	m := NewModel(apiURL)

	// Start Bubble Tea
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// Init initializes the model and triggers initial loading
func (m Model) Init() tea.Cmd {
	// tea.WithAltScreen() in main() already handles alt-screen; do not send
	// tea.EnterAltScreen here or Bubble Tea will render a blank screen first.
	return fetchOpenAPICmd(m.apiURL)
}

// Update handles model updates based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case clipboardNoticeExpiredMsg:
		m.clipboardMsg = ""
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// header = 2 rows, footer = 1 row
		headerHeight := 2
		footerHeight := 1
		availableHeight := m.height - headerHeight - footerHeight
		if availableHeight < 10 {
			availableHeight = 10
		}

		sidebarWidthPercent := 30
		sidebarWidth := (m.width * sidebarWidthPercent) / 100
		if sidebarWidth < 25 {
			sidebarWidth = 25
		}
		if sidebarWidth > m.width-20 {
			sidebarWidth = m.width - 20
		}

		rightWidth := m.width - sidebarWidth
		if rightWidth < 20 {
			rightWidth = 20
			sidebarWidth = m.width - rightWidth
		}

		editorHeightPercent := 40
		editorHeight := (availableHeight * editorHeightPercent) / 100
		if editorHeight < 5 {
			editorHeight = 5
		}

		responseHeight := availableHeight - editorHeight
		if responseHeight < 5 {
			responseHeight = 5
			editorHeight = availableHeight - responseHeight
		}

		// borderPadding = 2 border cols + 2 padding cols (Padding(0,1) each side)
		borderPadding := 4

		// ── Sidebar viewport: size it so handleNavigatingKeys has correct height ──
		sidebarVpW := sidebarWidth - borderPadding - 2 // -2 for Padding(0,1)
		if sidebarVpW < 1 {
			sidebarVpW = 1
		}
		sidebarVpH := availableHeight - 2 - 1 // -2 border, -1 title row
		if sidebarVpH < 1 {
			sidebarVpH = 1
		}
		m.sidebarView = viewport.New(sidebarVpW, sidebarVpH)

		// ── Response viewport: size it, then re-feed stored content ──
		responseInnerH := responseHeight - 2          // subtract border rows
		responseVpW := rightWidth - borderPadding - 2 // -2 for Padding(0,1)
		if responseVpW < 1 {
			responseVpW = 1
		}
		responseVpH := responseInnerH - 1 // -1 title row
		if responseVpH < 1 {
			responseVpH = 1
		}
		savedY := m.responseView.YOffset
		m.responseView = viewport.New(responseVpW, responseVpH)
		if m.responseContent != "" {
			m.responseView.SetContent(m.responseContent)
			m.responseView.SetYOffset(savedY)
		}

	case tea.KeyMsg:
		// Handle global keys
		switch msg.String() {
		case "ctrl+c", "q":
			if m.state == StateNavigating {
				return m, tea.Quit
			}
		}

		// Handle keys based on state
		switch m.state {
		case StateNavigating:
			return m.handleNavigatingKeys(msg)
		case StateEditingParams:
			return m.handleEditingKeys(msg)
		case StateViewingResponse:
			return m.handleViewingKeys(msg)
		case StateSearching:
			return m.handleSearchingKeys(msg)
		}

	case openAPISpecMsg:
		m.spec = msg.spec
		m.endpoints = ParseEndpoints(m.spec)
		m.baseURL = GetBaseURL(m.spec, m.apiURL)

		// Sort endpoints with the same ordering used by renderSidebarPanel so
		// j/k navigation matches the visual list exactly:
		//   1. Tag alphabetically
		//   2. Path alphabetically within tag
		//   3. Method order: GET → POST → PUT → DELETE → PATCH
		methodOrder := map[HTTPMethod]int{GET: 0, POST: 1, PUT: 2, DELETE: 3, PATCH: 4}
		sort.SliceStable(m.endpoints, func(i, j int) bool {
			ti := m.endpoints[i].Tag
			if ti == "" {
				ti = "default"
			}
			tj := m.endpoints[j].Tag
			if tj == "" {
				tj = "default"
			}
			if ti != tj {
				return ti < tj
			}
			if m.endpoints[i].Path != m.endpoints[j].Path {
				return m.endpoints[i].Path < m.endpoints[j].Path
			}
			return methodOrder[m.endpoints[i].Method] < methodOrder[m.endpoints[j].Method]
		})

		// Select first endpoint by default
		if len(m.endpoints) > 0 {
			m.selectedEndpoint = &m.endpoints[0]
			m = m.initializeInputs()
		}

		m.loading = false
		m.err = nil

	case errMsg:
		m.err = msg.err
		m.loading = false

	case apiResponseMsg:
		m.lastResponse = msg.response
		m.loading = false
		m.state = StateNavigating
		m.focusedPanel = 2

		// Build and store the highlighted content so it can be fed to
		// m.responseView now and also re-fed after every window resize.
		if msg.response != nil {
			styles := m.GetStyles()
			statusStyle := styles.StatusCodeStyle(msg.response.StatusCode)
			statusLine := fmt.Sprintf("Status: %s  Time: %s",
				statusStyle.Render(fmt.Sprintf("%d", msg.response.StatusCode)),
				styles.DescStyle().Render(fmt.Sprintf("%dms", msg.response.Duration.Milliseconds())),
			)
			var bodyStr string
			var jsonData interface{}
			if err := json.Unmarshal([]byte(msg.response.Body), &jsonData); err == nil {
				prettyBytes, _ := json.MarshalIndent(jsonData, "", "  ")
				bodyStr = highlightJSON(string(prettyBytes))
			} else {
				bodyStr = msg.response.Body
			}
			m.responseContent = statusLine + "\n\n" + bodyStr
		} else {
			m.responseContent = ""
		}
		m.responseView.SetContent(m.responseContent)
		m.responseView.GotoTop()

	case executeRequestMsg:
		m.loading = true
		body := m.buildRequestBody()
		return m, executeRequestCmd(m.baseURL, m.selectedEndpoint, body)
	}

	return m, tea.Batch(cmds...)
}

// filteredEndpoints returns the subset of m.endpoints that match the active
// searchQuery, or the full slice when no filter is active.
func (m Model) filteredEndpoints() []Endpoint {
	if m.searchQuery == "" {
		return m.endpoints
	}
	q := strings.ToLower(m.searchQuery)
	out := make([]Endpoint, 0)
	for _, ep := range m.endpoints {
		if strings.Contains(strings.ToLower(string(ep.Method)+" "+ep.Path+" "+ep.Tag+" "+ep.Summary), q) {
			out = append(out, ep)
		}
	}
	return out
}

// handleNavigatingKeys handles keys in navigation mode
func (m Model) handleNavigatingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// If focused on response panel, handle scroll keys via the stored viewport.
	if m.focusedPanel == 2 {
		switch msg.String() {
		case "j", "down":
			m.responseView.LineDown(1)
			return m, nil
		case "k", "up":
			m.responseView.LineUp(1)
			return m, nil
		case "d", "ctrl+d":
			m.responseView.HalfViewDown()
			return m, nil
		case "u", "ctrl+u":
			m.responseView.HalfViewUp()
			return m, nil
		case "g":
			m.responseView.GotoTop()
			return m, nil
		case "G":
			m.responseView.GotoBottom()
			return m, nil
		case "y":
			if m.lastResponse != nil {
				_ = clipboard.WriteAll(m.lastResponse.Body)
				m.clipboardMsg = " Copied!"
				return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
					return clipboardNoticeExpiredMsg{}
				})
			}
			return m, nil
		}
	}

	switch msg.String() {
	case "esc":
		// If a search filter is active, clear it and restore the full list
		if m.focusedPanel == 0 && m.searchQuery != "" {
			m.searchQuery = ""
			m.searchInput.SetValue("")
			return m, nil
		}

	case "j", "down":
		// Navigate down within the active filtered list
		if m.focusedPanel == 0 && m.selectedEndpoint != nil {
			list := m.filteredEndpoints()
			for i, ep := range list {
				if ep.Path == m.selectedEndpoint.Path && ep.Method == m.selectedEndpoint.Method {
					if i < len(list)-1 {
						next := list[i+1]
						// find the real index in m.endpoints for initializeInputs
						for ri := range m.endpoints {
							if m.endpoints[ri].Path == next.Path && m.endpoints[ri].Method == next.Method {
								m.listCursor = ri
								m.selectedEndpoint = &m.endpoints[ri]
								m = m.initializeInputs()
								break
							}
						}
					}
					break
				}
			}
		}

	case "k", "up":
		// Navigate up within the active filtered list
		if m.focusedPanel == 0 && m.selectedEndpoint != nil {
			list := m.filteredEndpoints()
			for i, ep := range list {
				if ep.Path == m.selectedEndpoint.Path && ep.Method == m.selectedEndpoint.Method {
					if i > 0 {
						prev := list[i-1]
						for ri := range m.endpoints {
							if m.endpoints[ri].Path == prev.Path && m.endpoints[ri].Method == prev.Method {
								m.listCursor = ri
								m.selectedEndpoint = &m.endpoints[ri]
								m = m.initializeInputs()
								break
							}
						}
					}
					break
				}
			}
		}

	case "enter":
		// Select endpoint and focus editor panel
		m.focusedPanel = 1
		m.state = StateNavigating

	case "r":
		// Execute request
		if m.selectedEndpoint != nil {
			return m, func() tea.Msg {
				return executeRequestMsg{}
			}
		}

	case "e":
		// Enter parameter editing mode
		if m.selectedEndpoint != nil && (len(m.selectedEndpoint.Parameters) > 0 || m.selectedEndpoint.RequestBody != nil) {
			m.state = StateEditingParams
			m.focusedPanel = 1
			m.currentInput = 0
			if len(m.paramInputs) > 0 {
				m.paramInputs[0].Focus()
			}
		}

	case "/":
		// Open sidebar search when focused on endpoint list
		if m.focusedPanel == 0 {
			// Compute sidebar vpWidth using the same layout math as renderMainArea
			sidebarWidth := (m.width * 30) / 100
			if sidebarWidth < 25 {
				sidebarWidth = 25
			}
			m.searchInput.Width = sidebarWidth - 6 // borderPadding(4) + input padding(2)
			m.state = StateSearching
			m.searchInput.SetValue("")
			m.searchQuery = ""
			m.searchInput.Focus()
			return m, nil
		}

	case "u":
		// Refresh the endpoint list by re-fetching the OpenAPI spec
		if m.focusedPanel == 0 {
			m.loading = true
			m.searchQuery = ""
			m.searchInput.SetValue("")
			return m, fetchOpenAPICmd(m.apiURL)
		}

	case "tab":
		// Change focused panel
		m.focusedPanel = (m.focusedPanel + 1) % 3

	case "shift+tab":
		// Change focused panel (reverse)
		m.focusedPanel = (m.focusedPanel - 1 + 3) % 3
	}

	return m, nil
}

// handleEditingKeys handles keys in parameter editing mode
func (m Model) handleEditingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "esc":
		// Exit editing mode and clear inputs
		m.state = StateNavigating
		if len(m.paramInputs) > 0 && m.currentInput < len(m.paramInputs) {
			m.paramInputs[m.currentInput].Blur()
		}
		// Clear all input fields
		for i := range m.paramInputs {
			m.paramInputs[i].SetValue("")
		}
		return m, nil

	case "ctrl+s":
		// Save and exit editing mode
		m = m.saveInputValues()
		m.state = StateNavigating
		if len(m.paramInputs) > 0 && m.currentInput < len(m.paramInputs) {
			m.paramInputs[m.currentInput].Blur()
		}
		return m, nil

	case "tab", "down", "enter":
		// Next input field
		if len(m.paramInputs) > 0 {
			m.paramInputs[m.currentInput].Blur()
			m.currentInput = (m.currentInput + 1) % len(m.paramInputs)
			m.paramInputs[m.currentInput].Focus()
		}
		return m, nil

	case "shift+tab", "up":
		// Previous input field
		if len(m.paramInputs) > 0 {
			m.paramInputs[m.currentInput].Blur()
			m.currentInput = (m.currentInput - 1 + len(m.paramInputs)) % len(m.paramInputs)
			m.paramInputs[m.currentInput].Focus()
		}
		return m, nil
	}

	// Update current input
	if len(m.paramInputs) > 0 && m.currentInput < len(m.paramInputs) {
		m.paramInputs[m.currentInput], cmd = m.paramInputs[m.currentInput].Update(msg)
	}
	return m, cmd
}

// handleViewingKeys handles keys in response viewing mode
func (m Model) handleViewingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.state = StateNavigating
		m.focusedPanel = 0
	case "r":
		if m.selectedEndpoint != nil {
			return m, func() tea.Msg { return executeRequestMsg{} }
		}
	case "j", "down":
		m.responseView.LineDown(1)
	case "k", "up":
		m.responseView.LineUp(1)
	case "d":
		m.responseView.HalfViewDown()
	case "u":
		m.responseView.HalfViewUp()
	case "g":
		m.responseView.GotoTop()
	case "G":
		m.responseView.GotoBottom()
	case "y":
		if m.lastResponse != nil {
			_ = clipboard.WriteAll(m.lastResponse.Body)
			m.clipboardMsg = " Copied!"
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return clipboardNoticeExpiredMsg{}
			})
		}
	}
	return m, nil
}

// View renders the model view
func (m Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	styles := m.GetStyles()

	// Render header and footer first so we can MEASURE their actual heights.
	// Using lipgloss.Height() is the only reliable way to know how many
	// terminal rows they consume — hardcoded constants drift when the style
	// gain extra borders or padding.
	header := m.renderHeader(styles)
	footer := m.renderFooter(styles)
	headerH := lipgloss.Height(header)
	footerH := lipgloss.Height(footer)

	availH := m.height - headerH - footerH
	if availH < 1 {
		availH = 1
	}

	var mainArea string
	if m.err != nil {
		mainArea = lipgloss.NewStyle().
			Width(m.width).
			Height(availH).MaxHeight(availH).
			Padding(1, 2).
			Render(
				styles.ErrorStyle().Render("❌ Error:\n\n") +
					m.err.Error() +
					"\n\n" +
					styles.DescStyle().Render("Press 'q' to quit"),
			)
	} else {
		mainArea = m.renderMainArea(availH, styles)
	}

	// Assemble the final frame and normalise to exactly m.width × m.height
	// cells.  The master container:
	//   • Height(m.height) — pads short content so we always output m.height rows
	//   • MaxHeight(m.height) — clips any overflow so we never exceed m.height rows
	//   • Background(colorSurface) — paints EVERY cell (incl. whitespace) with the
	//     dark background, which forces Bubble Tea’s diff renderer to repaint the
	//     entire screen each frame and eliminates ghost lines from stale renders.
	full := lipgloss.JoinVertical(lipgloss.Left, header, mainArea, footer)
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).MaxHeight(m.height).
		Background(lipgloss.Color(colorSurface)).
		Render(full)
}

// Custom messages

type openAPISpecMsg struct {
	spec *OpenAPISpec
}

type errMsg struct {
	err error
}

type apiResponseMsg struct {
	response *APIResponse
}

type executeRequestMsg struct{}

type clipboardNoticeExpiredMsg struct{}

// Commands

func fetchOpenAPICmd(apiURL string) tea.Cmd {
	return func() tea.Msg {
		spec, err := FetchOpenAPISpec(apiURL)
		if err != nil {
			return errMsg{err: err}
		}
		return openAPISpecMsg{spec: spec}
	}
}

func executeRequestCmd(baseURL string, endpoint *Endpoint, body string) tea.Cmd {
	return func() tea.Msg {
		response, _ := ExecuteRequest(baseURL, endpoint, body)
		return apiResponseMsg{response: response}
	}
}

// initTextInput initializes a text input for a parameter
func initTextInput(placeholder string) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 256
	ti.Width = 50
	return ti
}

// initViewport initializes the viewport for showing responses
func initViewport() viewport.Model {
	vp := viewport.New(60, 20)
	vp.SetContent("")
	return vp
}

// NewModel creates a new model with initialized components
func NewModel(apiURL string) Model {
	theme := ProfessionalTheme()
	si := textinput.New()
	si.Placeholder = "search endpoints..."
	si.CharLimit = 64

	return Model{
		apiURL:        apiURL,
		state:         StateNavigating,
		focusedPanel:  0,
		styles:        NewStyleConfig(theme),
		paramInputs:   []textinput.Model{},
		responseView:  initViewport(),
		sidebarView:   initViewport(),
		loading:       true,
		listCursor:    0,
		listOffset:    0,
		collapsedTags: make(map[string]bool),
		searchInput:   si,
	}
}

// handleSearchingKeys handles keypresses while the sidebar search box is active
func (m Model) handleSearchingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Cancel: clear query and return to navigation
		m.searchQuery = ""
		m.searchInput.SetValue("")
		m.searchInput.Blur()
		m.state = StateNavigating
		return m, nil
	case "enter":
		// Confirm: keep current query, apply first match as selection, exit search
		m.searchQuery = m.searchInput.Value()
		m.searchInput.Blur()
		m.state = StateNavigating
		// Move cursor to first visible endpoint that matches
		if m.searchQuery != "" {
			q := strings.ToLower(m.searchQuery)
			for i, ep := range m.endpoints {
				candidate := strings.ToLower(string(ep.Method) + " " + ep.Path + " " + ep.Tag + " " + ep.Summary)
				if strings.Contains(candidate, q) {
					m.listCursor = i
					m.selectedEndpoint = &m.endpoints[i]
					m = m.initializeInputs()
					break
				}
			}
		}
		return m, nil
	default:
		// Feed keystrokes to the textinput and keep query in sync
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		m.searchQuery = m.searchInput.Value()
		return m, cmd
	}
}

// initializeInputs creates input fields for the selected endpoint's parameters and body
func (m Model) initializeInputs() Model {
	if m.selectedEndpoint == nil {
		return m
	}

	m.paramInputs = []textinput.Model{}

	// Create inputs for parameters
	for _, param := range m.selectedEndpoint.Parameters {
		input := initTextInput(param.Name + " (" + param.Type + ")")
		if param.Value != "" {
			input.SetValue(param.Value)
		}
		m.paramInputs = append(m.paramInputs, input)
	}

	// Create inputs for body fields if request body exists
	if m.selectedEndpoint.RequestBody != nil {
		// Resolve schema if it's a $ref
		schema := resolveSchema(m.selectedEndpoint.RequestBody.Schema, m.spec)

		// Check if schema has properties
		if properties, ok := schema["properties"].(map[string]interface{}); ok {
			// Create inputs for each field
			var fieldNames []string
			for fieldName := range properties {
				fieldNames = append(fieldNames, fieldName)
			}
			sort.Strings(fieldNames) // Sort for consistent order

			for _, fieldName := range fieldNames {
				input := initTextInput(fieldName)
				m.paramInputs = append(m.paramInputs, input)
			}
		}
	}

	m.currentInput = 0
	return m
}

// saveInputValues saves input values back to parameters and request body
func (m Model) saveInputValues() Model {
	if m.selectedEndpoint == nil {
		return m
	}

	// Save regular parameter values
	for i := 0; i < len(m.selectedEndpoint.Parameters) && i < len(m.paramInputs); i++ {
		m.selectedEndpoint.Parameters[i].Value = m.paramInputs[i].Value()
	}

	return m
}

// resolveSchema resolves a schema, following $ref if present
func resolveSchema(schema map[string]interface{}, spec *OpenAPISpec) map[string]interface{} {
	if schema == nil || spec == nil {
		return schema
	}

	// Check for $ref
	if ref, ok := schema["$ref"].(string); ok {
		// Parse the reference (e.g., "#/components/schemas/Order")
		// Format: #/components/schemas/SchemaName
		if len(ref) > 21 && ref[:21] == "#/components/schemas/" {
			schemaName := ref[21:]
			if spec.Components != nil {
				if schemas, ok := spec.Components["schemas"]; ok {
					if resolvedSchema, ok := schemas[schemaName].(map[string]interface{}); ok {
						return resolvedSchema
					}
				}
			}
		}
	}

	return schema
}

// buildRequestBody builds the request body from input fields
func (m Model) buildRequestBody() string {
	if m.selectedEndpoint == nil || m.selectedEndpoint.RequestBody == nil {
		return ""
	}

	// Count how many inputs are for body fields (after parameters)
	numParams := len(m.selectedEndpoint.Parameters)
	if numParams >= len(m.paramInputs) {
		return "{}"
	}

	// Build JSON object from body field inputs
	body := make(map[string]interface{})

	// Resolve schema to handle $ref
	schema := resolveSchema(m.selectedEndpoint.RequestBody.Schema, m.spec)

	if properties, ok := schema["properties"].(map[string]interface{}); ok {
		// Iterate through the body field inputs (skip parameter inputs)
		for i := numParams; i < len(m.paramInputs); i++ {
			fieldName := m.paramInputs[i].Placeholder
			value := m.paramInputs[i].Value()

			if value != "" {
				// Try to get type from schema
				if fieldSchema, ok := properties[fieldName].(map[string]interface{}); ok {
					if fieldType, ok := fieldSchema["type"].(string); ok {
						switch fieldType {
						case "integer":
							var intVal int
							fmt.Sscanf(value, "%d", &intVal)
							body[fieldName] = intVal
						case "number":
							var floatVal float64
							fmt.Sscanf(value, "%f", &floatVal)
							body[fieldName] = floatVal
						case "boolean":
							body[fieldName] = value == "true"
						default:
							body[fieldName] = value
						}
					} else {
						body[fieldName] = value
					}
				} else {
					body[fieldName] = value
				}
			}
		}
	}

	jsonBytes, _ := json.Marshal(body)
	return string(jsonBytes)
}
