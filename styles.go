package main

import (
	"github.com/charmbracelet/lipgloss"
)

// Professional color palette inspired by Sway/LazyVim
var (
	// Dark mode base colors (The Void)
	colorBackground = "#0d1117"
	colorSurface     = "#0d1117" // Fondo base (el azul-negro profundo)
    colorBorder      = "#30363d" // Borde normal (sutil, estilo GitHub Dark)
    colorBorderDim   = "#161b22" // Borde casi invisible para paneles secundarios

    // Focus colors (The Glow) - Usando el Violeta para consistencia con POST
    colorFocusBorder = "#bb9af7" // Un violeta eléctrico que resalta más que el azul
    colorFocusGlow   = "#bb9af7" // El "brillo" que usaremos en títulos o indicadores

	// Accent colors
	colorPrimary   = "#7aa2f7" // Blue Focus (soft but visible)
    colorSecondary = "#bb9af7" // Violet accent (vivid)
    colorSuccess   = "#9ece6a" // GET Emerald
    colorError     = "#f7768e" // DELETE Coral/Red
    colorWarning   = "#e0af68" // PUT Amber/Gold
    colorInfo      = "#0db9d7" // Cyan info

    // Text colors (Hierarchy)
    colorTextPrimary   = "#c0caf5" // High readability
    colorTextSecondary = "#565f89" // Darker gray for prose/inactive
    colorTextMuted     = "#414868" // Very muted for borders/decorations
    colorTextHighlight = "#73daca" // Glow accent (teal)

    // Method colors (The "Glow" from Image 1)
    colorMethodGET    = "#73daca" // Cyan/Teal (from image)
    colorMethodPOST   = "#bb9af7" // Purple (from image)
    colorMethodPUT    = "#e0af68" // Orange/Amber (from image)
    colorMethodDELETE = "#f7768e" // Red/Pink (from image)
    colorMethodPATCH  = "#b4f9f8" // Soft Cyan/White (contrast)

    // Syntax highlighting (Optimized for #0d1117)
    colorJSONKey    = "#7aa2f7" // Soft Blue
    colorJSONString = "#9ece6a" // Green
    colorJSONNumber = "#ff9e64" // Orange
    colorJSONBool   = "#bb9af7" // Purple
    colorJSONNull   = "#565f89" // Muted Gray
    colorJSONPunct  = "#444b6a" // Subtle separation
)

// AppTheme defines the application theme
type AppTheme struct {
	// Base colors
	Background string
	Surface    string

	// Border colors
	Border      string
	BorderFocus string
	BorderDim   string

	// UI colors
	Primary   string
	Secondary string
	Success   string
	Error     string
	Warning   string
	Info      string

	// Text colors
	Foreground string
	Muted      string
	Highlight  string

	// Method colors
	MethodGET    string
	MethodPOST   string
	MethodPUT    string
	MethodDELETE string
	MethodPATCH  string

	// Syntax colors
	JSONKey    string
	JSONString string
	JSONNumber string
	JSONBool   string
	JSONNull   string
	JSONPunct  string
}

// ProfessionalTheme returns the high-end dark theme
func ProfessionalTheme() AppTheme {
	return AppTheme{
		Background:   colorBackground,
		Surface:      colorSurface,
		Border:       colorBorder,
		BorderFocus:  colorFocusBorder,
		BorderDim:    colorBorderDim,
		Primary:      colorPrimary,
		Secondary:    colorSecondary,
		Success:      colorSuccess,
		Error:        colorError,
		Warning:      colorWarning,
		Info:         colorInfo,
		Foreground:   colorTextPrimary,
		Muted:        colorTextMuted,
		Highlight:    colorTextHighlight,
		MethodGET:    colorMethodGET,
		MethodPOST:   colorMethodPOST,
		MethodPUT:    colorMethodPUT,
		MethodDELETE: colorMethodDELETE,
		MethodPATCH:  colorMethodPATCH,
		JSONKey:      colorJSONKey,
		JSONString:   colorJSONString,
		JSONNumber:   colorJSONNumber,
		JSONBool:     colorJSONBool,
		JSONNull:     colorJSONNull,
		JSONPunct:    colorJSONPunct,
	}
}

// StyleConfig contains all lipgloss styles
type StyleConfig struct {
	theme AppTheme

	// Panel styles
	PanelBase      lipgloss.Style
	PanelFocused   lipgloss.Style
	PanelUnfocused lipgloss.Style

	// Title styles
	TitleFocused   lipgloss.Style
	TitleUnfocused lipgloss.Style

	// Method label styles
	MethodGET    lipgloss.Style
	MethodPOST   lipgloss.Style
	MethodPUT    lipgloss.Style
	MethodDELETE lipgloss.Style
	MethodPATCH  lipgloss.Style

	// Text styles
	TextNormal    lipgloss.Style
	TextMuted     lipgloss.Style
	TextHighlight lipgloss.Style
	TextSuccess   lipgloss.Style
	TextError     lipgloss.Style
	TextWarning   lipgloss.Style
	TextInfo      lipgloss.Style

	// Input styles
	InputBox        lipgloss.Style
	InputBoxFocused lipgloss.Style
	InputLabel      lipgloss.Style

	// Status bar
	StatusBar      lipgloss.Style
	StatusBarKey   lipgloss.Style
	StatusBarValue lipgloss.Style
	StatusBarSep   lipgloss.Style

	// Endpoint list
	EndpointTag      lipgloss.Style
	EndpointSelected lipgloss.Style
	EndpointNormal   lipgloss.Style
}

// NewStyleConfig creates a new style configuration
func NewStyleConfig(theme AppTheme) StyleConfig {
	s := StyleConfig{theme: theme}

	// Base panel style
	s.PanelBase = lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder())

	// Focused panel
	s.PanelFocused = s.PanelBase.Copy().
		BorderForeground(lipgloss.Color(theme.BorderFocus))

	// Unfocused panel
	s.PanelUnfocused = s.PanelBase.Copy().
		BorderForeground(lipgloss.Color(theme.Border))

	// Titles
	s.TitleFocused = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.Primary)).
		Background(lipgloss.Color(theme.Surface)).
		Padding(0, 1)

	s.TitleUnfocused = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.Muted)).
		Padding(0, 1)

	// Method labels with rounded style
	methodLabelBase := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	s.MethodGET = methodLabelBase.Copy().
		Foreground(lipgloss.Color(theme.MethodGET))

	s.MethodPOST = methodLabelBase.Copy().
		Foreground(lipgloss.Color(theme.MethodPOST))

	s.MethodPUT = methodLabelBase.Copy().
		Foreground(lipgloss.Color(theme.MethodPUT))

	s.MethodDELETE = methodLabelBase.Copy().
		Foreground(lipgloss.Color(theme.MethodDELETE))

	s.MethodPATCH = methodLabelBase.Copy().
		Foreground(lipgloss.Color(theme.MethodPATCH))

	// Text styles
	s.TextNormal = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground))

	s.TextMuted = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Muted))

	s.TextHighlight = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Highlight)).
		Bold(true)

	s.TextSuccess = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Success))

	s.TextError = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Error))

	s.TextWarning = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Warning))

	s.TextInfo = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Info))

	// Input box styles
	s.InputBox = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderDim)).
		Padding(0, 1)

	s.InputBoxFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.BorderFocus)).
		Padding(0, 1)

	s.InputLabel = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Secondary)).
		Bold(true).
		MarginBottom(1)

	// Status bar
	s.StatusBar = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground)).
		Background(lipgloss.Color(theme.Surface)).
		Padding(0, 1)

	s.StatusBarKey = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Info)).
		Bold(true)

	s.StatusBarValue = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground))

	s.StatusBarSep = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.BorderDim))

	// Endpoint list
	s.EndpointTag = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Highlight)).
		Bold(true).
		MarginTop(1).
		MarginBottom(1)

	s.EndpointSelected = lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Surface)).
		Foreground(lipgloss.Color(theme.Foreground)).
		Bold(true)

	s.EndpointNormal = lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Foreground))

	return s
}

// GetMethodStyle returns the appropriate style for an HTTP method
func (s StyleConfig) GetMethodStyle(method string) lipgloss.Style {
	switch method {
	case "GET":
		return s.MethodGET
	case "POST":
		return s.MethodPOST
	case "PUT":
		return s.MethodPUT
	case "DELETE":
		return s.MethodDELETE
	case "PATCH":
		return s.MethodPATCH
	default:
		return s.TextNormal
	}
}

// GetStatusCodeStyle returns the appropriate style for a status code
func (s StyleConfig) GetStatusCodeStyle(code int) lipgloss.Style {
	switch {
	case code >= 200 && code < 300:
		return s.TextSuccess
	case code >= 300 && code < 400:
		return s.TextInfo
	case code >= 400 && code < 500:
		return s.TextWarning
	case code >= 500:
		return s.TextError
	default:
		return s.TextMuted
	}
}
