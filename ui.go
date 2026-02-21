package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// resolveSchemaUI resolves a schema, following $ref if present (UI version)
func resolveSchemaUI(schema map[string]interface{}, spec *OpenAPISpec) map[string]interface{} {
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

// Styles is a backwards-compatible wrapper around StyleConfig
type Styles struct {
	BorderColor string
	Primary     string
	Secondary   string
	Success     string
	Error       string
	Warning     string
	Info        string
	Foreground  string
	Muted       string
	Highlight   string
}

// GetStyles returns a simplified styles struct for backwards compatibility
func (m Model) GetStyles() Styles {
	return Styles{
		BorderColor: m.styles.theme.Border,
		Primary:     m.styles.theme.Primary,
		Secondary:   m.styles.theme.Secondary,
		Success:     m.styles.theme.Success,
		Error:       m.styles.theme.Error,
		Warning:     m.styles.theme.Warning,
		Info:        m.styles.theme.Info,
		Foreground:  m.styles.theme.Foreground,
		Muted:       m.styles.theme.Muted,
		Highlight:   m.styles.theme.Highlight,
	}
}

// PanelStyle returns the base style for a panel with focus indication.
// Layout contract:
//   - .Width(w)  → lipgloss wraps text at w cols; outer rendered = w + 2 padding + 2 border = w+4 cols
//   - .Height(h) → pads content to h rows;      outer rendered = h + 2 border rows
//   - MaxWidth / MaxHeight hard-clips so no panel can ever push siblings
func (s Styles) PanelStyle(focused bool) lipgloss.Style {
	borderColor := s.BorderColor
	if focused {
		borderColor = s.Primary
	}
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1)
}

// TitleStyle returns the style for panel titles
func (s Styles) TitleStyle(focused bool) lipgloss.Style {
	color := s.Muted
	if focused {
		color = s.Primary
	}
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Padding(0, 1)
}

// MethodStyle returns the style for an HTTP method
func (s Styles) MethodStyle(method HTTPMethod) lipgloss.Style {
	var color string
	switch method {
	case GET:
		color = "#3fb950"
	case POST:
		color = "#bc8cff"
	case PUT:
		color = "#d29922"
	case DELETE:
		color = "#f85149"
	case PATCH:
		color = "#58a6ff"
	default:
		color = s.Muted
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true).
		Width(7).
		Align(lipgloss.Center)
}

// StatusCodeStyle returns the style for an HTTP status code
func (s Styles) StatusCodeStyle(statusCode int) lipgloss.Style {
	var color string
	switch {
	case statusCode >= 200 && statusCode < 300:
		color = s.Success
	case statusCode >= 300 && statusCode < 400:
		color = s.Info
	case statusCode >= 400 && statusCode < 500:
		color = s.Warning
	case statusCode >= 500:
		color = s.Error
	default:
		color = s.Muted
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Bold(true)
}

// KeyStyle returns the style for help keys
func (s Styles) KeyStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Highlight)).
		Bold(true)
}

// DescStyle returns the style for help descriptions
func (s Styles) DescStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Muted))
}

// ErrorStyle returns the style for error messages
func (s Styles) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Error)).
		Bold(true)
}

// TagStyle returns the style for endpoint tags
func (s Styles) TagStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(s.Highlight)).
		Bold(true).
		Italic(true)
}

// renderMainArea renders the three-panel grid into exactly availH rows.
// The caller (View) is responsible for header and footer assembly so that
// measured heights are used — no hardcoded offsets that can drift.
func (m Model) renderMainArea(availH int, styles Styles) string {
	if availH < 10 {
		availH = 10
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

	borderPadding := 4
	sidebarInnerWidth := sidebarWidth - borderPadding
	rightInnerWidth := rightWidth - borderPadding
	if sidebarInnerWidth < 1 {
		sidebarInnerWidth = 1
	}
	if rightInnerWidth < 1 {
		rightInnerWidth = 1
	}

	editorHeightPercent := 40
	editorHeight := (availH * editorHeightPercent) / 100
	if editorHeight < 5 {
		editorHeight = 5
	}

	responseHeight := availH - editorHeight
	if responseHeight < 5 {
		responseHeight = 5
		editorHeight = availH - responseHeight
	}

	editorInnerHeight := editorHeight - 2
	responseInnerHeight := responseHeight - 2
	sidebarInnerHeight := availH - 2

	if editorInnerHeight < 1 {
		editorInnerHeight = 1
	}
	if responseInnerHeight < 1 {
		responseInnerHeight = 1
	}
	if sidebarInnerHeight < 1 {
		sidebarInnerHeight = 1
	}

	sidebarPanel := m.renderSidebarPanel(sidebarInnerWidth, sidebarInnerHeight, styles)
	editorPanel := m.renderEditorPanel(rightInnerWidth, editorInnerHeight, styles)
	responsePanel := m.renderResponsePanel(rightInnerWidth, responseInnerHeight, styles)

	// Stack the two right-side panels, then place the sidebar to their left.
	rightColumn := lipgloss.JoinVertical(lipgloss.Top, editorPanel, responsePanel)
	return lipgloss.JoinHorizontal(lipgloss.Top, sidebarPanel, rightColumn)
}

// renderHeader renders the application header
func (m Model) renderHeader(styles Styles) string {
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.Primary)).
		Bold(true).
		Padding(0, 2)

	infoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.Muted)).
		Italic(true)

	title := titleStyle.Render("⚡ SNAG")
	info := ""

	if m.spec != nil {
		info = infoStyle.Render(fmt.Sprintf("  %s v%s", m.spec.Info.Title, m.spec.Info.Version))
	}

	header := lipgloss.JoinHorizontal(lipgloss.Top, title, info)

	return lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 2).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(styles.BorderColor)).
		Render(header)
}

// renderFooter renders the footer with keyboard help (nvim-style with status)
func (m Model) renderFooter(styles Styles) string {
	keys := []string{}

	switch m.state {
	case StateSearching:
		keys = []string{
			styles.KeyStyle().Render("type") + styles.DescStyle().Render(" filter"),
			styles.KeyStyle().Render("enter") + styles.DescStyle().Render(" confirm"),
			styles.KeyStyle().Render("esc") + styles.DescStyle().Render(" cancel"),
		}
	case StateNavigating:
		if m.focusedPanel == 2 {
			// When focused on Response panel, show scroll keys
			keys = []string{
				styles.KeyStyle().Render("j/k") + styles.DescStyle().Render(" scroll"),
				styles.KeyStyle().Render("d/u") + styles.DescStyle().Render(" page"),
				styles.KeyStyle().Render("g/G") + styles.DescStyle().Render(" top/bottom"),
				styles.KeyStyle().Render("y") + styles.DescStyle().Render(" copy"),
				styles.KeyStyle().Render("tab") + styles.DescStyle().Render(" switch panel"),
			}
		} else if m.focusedPanel == 0 && m.searchQuery != "" {
			keys = []string{
				styles.KeyStyle().Render("j/k") + styles.DescStyle().Render(" navigate"),
				styles.KeyStyle().Render("/") + styles.DescStyle().Render(" search"),
				styles.KeyStyle().Render("esc") + styles.DescStyle().Render(" clear filter"),
				styles.KeyStyle().Render("r") + styles.DescStyle().Render(" run"),
			}
		} else {
			keys = []string{
				styles.KeyStyle().Render("j/k") + styles.DescStyle().Render(" navigate"),
				styles.KeyStyle().Render("/") + styles.DescStyle().Render(" search"),
				styles.KeyStyle().Render("tab") + styles.DescStyle().Render(" switch panel"),
				styles.KeyStyle().Render("r") + styles.DescStyle().Render(" run"),
				styles.KeyStyle().Render("e") + styles.DescStyle().Render(" edit"),
			}
		}
	case StateEditingParams:
		keys = []string{
			styles.KeyStyle().Render("tab") + styles.DescStyle().Render(" next field"),
			styles.KeyStyle().Render("shift+tab") + styles.DescStyle().Render(" prev field"),
			styles.KeyStyle().Render("ctrl+s") + styles.DescStyle().Render(" save"),
			styles.KeyStyle().Render("esc") + styles.DescStyle().Render(" cancel"),
		}
	case StateViewingResponse:
		keys = []string{
			styles.KeyStyle().Render("j/k") + styles.DescStyle().Render(" scroll"),
			styles.KeyStyle().Render("y") + styles.DescStyle().Render(" copy"),
			styles.KeyStyle().Render("esc") + styles.DescStyle().Render(" back"),
			styles.KeyStyle().Render("r") + styles.DescStyle().Render(" re-run"),
		}
	}

	// Add status information on the right side
	status := ""
	if m.clipboardMsg != "" {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Success)).Render(m.clipboardMsg)
	} else if m.lastResponse != nil {
		statusCode := styles.StatusCodeStyle(m.lastResponse.StatusCode).Render(fmt.Sprintf("%d", m.lastResponse.StatusCode))
		duration := styles.DescStyle().Render(fmt.Sprintf(" %dms", m.lastResponse.Duration.Milliseconds()))
		status = " " + statusCode + duration
	}

	// Join with nvim-style separators
	sep := lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BorderColor)).Render("  ·  ")
	help := strings.Join(keys, sep) + status

	return lipgloss.NewStyle().
		Width(m.width).
		Padding(0, 2).
		Foreground(lipgloss.Color(styles.Muted)).
		Render(help)
}

// renderSidebarPanel renders the endpoint navigation panel with collapsible groups
func (m Model) renderSidebarPanel(width, innerHeight int, styles Styles) string {
	focused := m.focusedPanel == 0
	// MaxWidth/MaxHeight ensure the panel never overflows its allocated grid cell.
	// Outer rendered size = width+4 cols × innerHeight+2 rows.
	panelStyle := styles.PanelStyle(focused).
		Width(width).MaxWidth(width + 4).
		Height(innerHeight).MaxHeight(innerHeight + 2)

	titleText := "📋 Endpoints"
	visibleEndpoints := m.filteredEndpoints()
	if m.searchQuery != "" {
		titleText = fmt.Sprintf("📋 Endpoints (%d)", len(visibleEndpoints))
	}

	content := ""
	cursorLine := 0 // actual rendered line of the selected endpoint (used for scroll)
	if m.loading && len(m.endpoints) == 0 {
		content = styles.DescStyle().Render("Loading endpoints...")
	} else if len(visibleEndpoints) == 0 {
		if m.searchQuery != "" {
			content = styles.DescStyle().Render("No matches for \"" + m.searchQuery + "\"")
		} else {
			content = styles.ErrorStyle().Render("No endpoints found")
		}
	} else {
		groups := make(map[string][]Endpoint)
		tagOrder := []string{}
		seen := make(map[string]bool)

		for _, ep := range visibleEndpoints {
			tag := ep.Tag
			if tag == "" {
				tag = "default"
			}
			groups[tag] = append(groups[tag], ep)
			if !seen[tag] {
				tagOrder = append(tagOrder, tag)
				seen[tag] = true
			}
		}

		methodOrder := map[HTTPMethod]int{
			GET:    0,
			POST:   1,
			PUT:    2,
			DELETE: 3,
			PATCH:  4,
		}
		sort.Strings(tagOrder)
		for _, tag := range tagOrder {
			endpoints := groups[tag]
			sort.SliceStable(endpoints, func(i, j int) bool {
				if endpoints[i].Path == endpoints[j].Path {
					oi := methodOrder[endpoints[i].Method]
					oj := methodOrder[endpoints[j].Method]
					return oi < oj
				}
				return endpoints[i].Path < endpoints[j].Path
			})
			groups[tag] = endpoints
		}

		var lines []string
		lineIdx := 0 // running rendered line counter
		for _, tag := range tagOrder {
			tagStyle := styles.TagStyle()
			collapsed, isCollapsed := m.collapsedTags[tag]
			icon := "▼"
			if isCollapsed && collapsed {
				icon = "▶"
			}
			lines = append(lines, tagStyle.Render(fmt.Sprintf("%s %s", icon, strings.ToUpper(tag))))
			lineIdx++ // tag header occupies one rendered line

			if m.collapsedTags[tag] {
				continue
			}

			for _, ep := range groups[tag] {
				isSelected := m.selectedEndpoint != nil &&
					ep.Path == m.selectedEndpoint.Path &&
					ep.Method == m.selectedEndpoint.Method

				if isSelected {
					cursorLine = lineIdx
				}

				methodStyle := styles.MethodStyle(ep.Method)
				methodLabel := methodStyle.Render(string(ep.Method))

				pathText := ep.Path
				maxPathLen := width - 10
				if maxPathLen < 5 {
					maxPathLen = 5
				}
				if len(pathText) > maxPathLen {
					pathText = pathText[:maxPathLen-3] + "..."
				}

				line := fmt.Sprintf("  %s %s", methodLabel, pathText)
				if isSelected {
					line = lipgloss.NewStyle().
						Foreground(lipgloss.Color(styles.Foreground)).
						Bold(true).
						Render(fmt.Sprintf("> %s %s", methodLabel, pathText))
				}
				lines = append(lines, line)
				lineIdx++
			}
		}

		content = strings.Join(lines, "\n")
	}

	// vpWidth = width - 2: the Padding(0,1) inside PanelStyle consumes 1 col on
	// each side, leaving width-2 cols for actual text content. Feeding the
	// viewport more than this causes lipgloss to word-wrap the viewport output
	// inside the panel, adding phantom extra lines and breaking the layout.
	vpWidth := width - 2
	if vpWidth < 1 {
		vpWidth = 1
	}
	// When the search box is visible it occupies 1 extra row below the title.
	searchVisible := m.state == StateSearching || m.searchQuery != ""
	titleRows := 1
	if searchVisible {
		titleRows = 2
	}
	vpHeight := innerHeight - titleRows
	if vpHeight < 1 {
		vpHeight = 1
	}

	vp := viewport.New(vpWidth, vpHeight)
	vp.SetContent(content)
	// Derive scroll offset from the actual rendered cursor line so that
	// the selected endpoint is always visible, regardless of tag headers.
	scrollY := m.listOffset
	if cursorLine < scrollY {
		scrollY = cursorLine
	} else if cursorLine >= scrollY+vpHeight {
		scrollY = cursorLine - vpHeight + 1
	}
	if scrollY < 0 {
		scrollY = 0
	}
	vp.SetYOffset(scrollY)

	title := styles.TitleStyle(focused).Render(titleText)
	header := title
	if searchVisible {
		// Style the search input to fit the panel width
		searchBar := lipgloss.NewStyle().
			Width(vpWidth).
			Foreground(lipgloss.Color(styles.Foreground)).
			Render(m.searchInput.View())
		header = title + "\n" + searchBar
	}
	return panelStyle.Render(header + "\n" + vp.View())
}

// renderEditorPanel renders the parameter/body editor panel with ghost inputs
func (m Model) renderEditorPanel(width, innerHeight int, styles Styles) string {
	focused := m.focusedPanel == 1
	panelStyle := styles.PanelStyle(focused).
		Width(width).MaxWidth(width + 4).
		Height(innerHeight).MaxHeight(innerHeight + 2)

	title := styles.TitleStyle(focused).Render("✏️  Request")

	// inputWidgetLine[i] records which line index (0-based) the actual text-input
	// widget for paramInputs[i] lands on inside the content string.  Used below
	// to compute the viewport scroll offset that keeps the focused field visible.
	inputWidgetLine := make([]int, 0, 8)
	lineCounter := 0 // counts lines as we append to info

	content := ""
	if m.selectedEndpoint == nil {
		content = styles.DescStyle().Render("Select an endpoint to begin")
	} else {
		info := []string{
			styles.KeyStyle().Render(string(m.selectedEndpoint.Method)) + " " + m.selectedEndpoint.Path,
		}
		lineCounter++ // method+path line

		if m.selectedEndpoint.Summary != "" {
			summary := m.selectedEndpoint.Summary
			maxLen := width - 4
			if maxLen < 10 {
				maxLen = 10
			}
			if len(summary) > maxLen {
				summary = summary[:maxLen-3] + "..."
			}
			info = append(info, styles.DescStyle().Render(summary))
			lineCounter++
		}

		if len(m.selectedEndpoint.Parameters) > 0 || m.selectedEndpoint.RequestBody != nil {
			info = append(info, "")
			lineCounter++
			info = append(info, styles.TagStyle().Render("Parameters:"))
			lineCounter++

			if m.state == StateEditingParams {
				for i, input := range m.paramInputs {
					var paramName, paramType string
					var required bool

					if i < len(m.selectedEndpoint.Parameters) {
						param := m.selectedEndpoint.Parameters[i]
						paramName = param.Name
						paramType = param.Type
						required = param.Required
					} else {
						paramName = input.Placeholder
						paramType = "string"

						if m.selectedEndpoint.RequestBody != nil {
							schema := resolveSchemaUI(m.selectedEndpoint.RequestBody.Schema, m.spec)
							if properties, ok := schema["properties"].(map[string]interface{}); ok {
								if fieldSchema, ok := properties[paramName].(map[string]interface{}); ok {
									if fieldType, ok := fieldSchema["type"].(string); ok {
										paramType = fieldType
									}
								}
							}
						}
					}

					label := paramName
					if required {
						label += styles.ErrorStyle().Render(" *")
					}
					label += styles.DescStyle().Render(" [" + paramType + "]")

					isActive := focused && i == m.currentInput
					cursorColor := styles.Primary
					textColor := styles.Foreground
					placeholderColor := styles.Muted
					if isActive {
						textColor = styles.Primary
					}

					ti := input
					ti.Width = width - 8
					if ti.Width < 5 {
						ti.Width = 5
					}
					ti.Prompt = ""
					ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Muted))
					ti.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cursorColor))
					ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(textColor))
					ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(placeholderColor))

					inputView := ti.View()
					info = append(info, "  "+label)
					lineCounter++ // label line
					info = append(info, "    "+inputView)
					inputWidgetLine = append(inputWidgetLine, lineCounter)
					lineCounter++ // widget line
				}
			} else {
				for i, param := range m.selectedEndpoint.Parameters {
					val := param.Value
					if i < len(m.paramInputs) && m.paramInputs[i].Value() != "" {
						val = m.paramInputs[i].Value()
					}
					paramLine := fmt.Sprintf("  %s %s",
						lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Muted)).Render(param.Name+" ["+param.Type+"]"),
						"",
					)
					if param.Required {
						paramLine = "  " + lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Muted)).Render(param.Name+" ["+param.Type+"]") + styles.ErrorStyle().Render(" *")
					}
					if val != "" {
						valDisplay := val
						maxValLen := width - len(param.Name) - 15
						if maxValLen < 5 {
							maxValLen = 5
						}
						if len(valDisplay) > maxValLen {
							valDisplay = valDisplay[:maxValLen-3] + "..."
						}
						paramLine += "  " + styles.KeyStyle().Render(valDisplay)
					} else {
						paramLine += "  " + styles.DescStyle().Render("(empty)")
					}
					info = append(info, paramLine)
				}

				if m.selectedEndpoint.RequestBody != nil {
					info = append(info, "")
					info = append(info, styles.TagStyle().Render("Request Body:"))
					numParams := len(m.selectedEndpoint.Parameters)
					hasBodyValues := false
					for i := numParams; i < len(m.paramInputs); i++ {
						v := m.paramInputs[i].Value()
						name := m.paramInputs[i].Placeholder
						if v != "" {
							valDisplay := v
							maxValLen := width - len(name) - 10
							if maxValLen < 5 {
								maxValLen = 5
							}
							if len(valDisplay) > maxValLen {
								valDisplay = valDisplay[:maxValLen-3] + "..."
							}
							info = append(info, "  "+lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Muted)).Render(name+":")+"  "+styles.KeyStyle().Render(valDisplay))
							hasBodyValues = true
						} else {
							info = append(info, "  "+lipgloss.NewStyle().Foreground(lipgloss.Color(styles.Muted)).Render(name+":")+"  "+styles.DescStyle().Render("(empty)"))
						}
					}
					if !hasBodyValues && len(m.paramInputs) <= numParams {
						info = append(info, styles.DescStyle().Render("  Press 'e' to edit"))
					}
				}
			}
		} else {
			info = append(info, "")
			info = append(info, styles.DescStyle().Render("No parameters required"))
		}

		content = strings.Join(info, "\n")
	}

	// vpWidth = width-2: accounts for Padding(0,1) inside PanelStyle
	vpWidth := width - 2
	if vpWidth < 1 {
		vpWidth = 1
	}
	vpHeight := innerHeight - 1
	if vpHeight < 1 {
		vpHeight = 1
	}

	vp := viewport.New(vpWidth, vpHeight)
	vp.SetContent(content)

	// Auto-scroll so the focused input field is always visible.
	// Each input occupies 2 lines: the label (widget-1) and the widget itself.
	// We scroll just enough so both lines fit inside the viewport window.
	editorScrollY := 0
	if m.state == StateEditingParams && m.currentInput < len(inputWidgetLine) {
		widgetLine := inputWidgetLine[m.currentInput] // 0-based line of the widget
		labelLine := widgetLine - 1                   // label is one line above

		// If the widget line is below the visible window, scroll down to it.
		if widgetLine >= editorScrollY+vpHeight {
			editorScrollY = widgetLine - vpHeight + 1
		}
		// If the label line is above the visible window, scroll up to it.
		if labelLine < editorScrollY {
			editorScrollY = labelLine
		}
	}
	vp.SetYOffset(editorScrollY)

	return panelStyle.Render(title + "\n" + vp.View())
}

// renderResponsePanel renders the API response panel with syntax highlighting
func (m Model) renderResponsePanel(width, innerHeight int, styles Styles) string {
	focused := m.focusedPanel == 2
	panelStyle := styles.PanelStyle(focused).
		Width(width).MaxWidth(width + 4).
		Height(innerHeight).MaxHeight(innerHeight + 2)

	spinner := ""
	if m.loading {
		spinner = styles.DescStyle().Render(" ⏳")
	}

	// Build a fresh viewport every frame, sized from the parameters that
	// renderMainArea computed for THIS frame.  This is the only guarantee
	// that View() output is exactly vpHeight rows — stale stored-viewport
	// dimensions diverge when the terminal is resized mid-scroll.
	//
	// vpWidth = width-2: Padding(0,1) inside PanelStyle consumes 1 col each side.
	vpWidth := width - 2
	if vpWidth < 1 {
		vpWidth = 1
	}
	vpHeight := innerHeight - 1 // 1 row for the title line
	if vpHeight < 1 {
		vpHeight = 1
	}

	vp := viewport.New(vpWidth, vpHeight)

	var content string
	if m.responseContent != "" {
		content = m.responseContent
	} else if !m.loading {
		content = styles.DescStyle().Render("Press 'r' to execute the request")
	}
	vp.SetContent(content)

	// Restore scroll position from the stored viewport (updated by key handlers).
	// SetYOffset auto-clamps to [0, TotalLineCount-vpHeight].
	vp.SetYOffset(m.responseView.YOffset)

	// Scroll-percentage hint — only when content overflows.
	scrollHint := ""
	if vp.TotalLineCount() > vpHeight {
		pct := int(vp.ScrollPercent() * 100)
		scrollHint = " " + styles.DescStyle().Render(fmt.Sprintf("%d%%", pct))
	}

	title := styles.TitleStyle(focused).Render("📡 Response") + spinner + scrollHint
	return panelStyle.Render(title + "\n" + vp.View())
}

// tokenKind identifies the JSON token type for highlighting.
type tokenKind int

const (
	kindNone   tokenKind = iota
	kindKey              // object key string
	kindString           // non-key string value
	kindNumber           // numeric literal
	kindBool             // true / false
	kindNull             // null
	kindPunct            // { } [ ] : ,
	kindSpace            // whitespace (passed through raw, no ANSI)
)

// highlightJSON applies syntax highlighting to JSON.
//
// Critical design rule: accumulate entire tokens into a single strings.Builder
// before calling lipgloss.Render ONCE per token.  The old character-by-character
// approach produced O(n) ANSI escape sequences for n input bytes, which caused
// the bubbles/viewport line-width measurement to diverge for large payloads
// (GET /users with 100+ objects) and made the panel overflow its allocated rows.
func highlightJSON(jsonStr string) string {
	if jsonStr == "" {
		return ""
	}

	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONKey))
	stringStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONString))
	numberStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONNumber))
	boolStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONBool))
	nullStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONNull))
	punctStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorJSONPunct))

	applyStyle := func(kind tokenKind, tok string) string {
		switch kind {
		case kindKey:
			return keyStyle.Render(tok)
		case kindString:
			return stringStyle.Render(tok)
		case kindNumber:
			return numberStyle.Render(tok)
		case kindBool:
			return boolStyle.Render(tok)
		case kindNull:
			return nullStyle.Render(tok)
		case kindPunct:
			return punctStyle.Render(tok)
		default:
			return tok
		}
	}

	var out strings.Builder
	out.Grow(len(jsonStr) * 2)

	i := 0
	n := len(jsonStr)

	for i < n {
		ch := jsonStr[i]

		// ── whitespace / newlines: emit raw (no ANSI, preserves line structure) ──
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			out.WriteByte(ch)
			i++
			continue
		}

		// ── punctuation: single-char tokens ──
		if ch == '{' || ch == '}' || ch == '[' || ch == ']' || ch == ':' || ch == ',' {
			out.WriteString(applyStyle(kindPunct, string(ch)))
			i++
			continue
		}

		// ── string (key or value): scan to closing unescaped quote ──
		if ch == '"' {
			start := i
			i++ // skip opening quote
			for i < n {
				if jsonStr[i] == '\\' {
					i += 2 // skip escaped character
					continue
				}
				if jsonStr[i] == '"' {
					i++ // include closing quote
					break
				}
				i++
			}
			tok := jsonStr[start:i]

			// Determine if this is a key: skip whitespace and check for ':'
			kind := kindString
			j := i
			for j < n && (jsonStr[j] == ' ' || jsonStr[j] == '\t') {
				j++
			}
			if j < n && jsonStr[j] == ':' {
				kind = kindKey
			}
			out.WriteString(applyStyle(kind, tok))
			continue
		}

		// ── number: [-][digits][.digits][e[+-]digits] ──
		if ch >= '0' && ch <= '9' || ch == '-' {
			start := i
			for i < n {
				c := jsonStr[i]
				if c >= '0' && c <= '9' || c == '.' || c == '-' || c == '+' || c == 'e' || c == 'E' {
					i++
				} else {
					break
				}
			}
			out.WriteString(applyStyle(kindNumber, jsonStr[start:i]))
			continue
		}

		// ── keywords: true, false, null ──
		if i+4 <= n && jsonStr[i:i+4] == "null" {
			out.WriteString(applyStyle(kindNull, "null"))
			i += 4
			continue
		}
		if i+4 <= n && jsonStr[i:i+4] == "true" {
			out.WriteString(applyStyle(kindBool, "true"))
			i += 4
			continue
		}
		if i+5 <= n && jsonStr[i:i+5] == "false" {
			out.WriteString(applyStyle(kindBool, "false"))
			i += 5
			continue
		}

		// ── fallback: emit raw byte ──
		out.WriteByte(ch)
		i++
	}

	return out.String()
}

// SpinnerFrames are the loading spinner frames
var SpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
