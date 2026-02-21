package main

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
)

// AppState represents the application state
type AppState int

const (
	StateNavigating      AppState = iota // Navigating through endpoints
	StateEditingParams                   // Editing request parameters/body
	StateLoading                         // Loading response
	StateViewingResponse                 // Viewing response
	StateSearching                       // Typing in sidebar search box
)

// HTTPMethod represents HTTP methods
type HTTPMethod string

const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
)

// Endpoint represents an API endpoint
type Endpoint struct {
	Method      HTTPMethod
	Path        string
	Tag         string
	Summary     string
	Description string
	Parameters  []Parameter
	RequestBody *RequestBody
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string
	In          string // query, path, header
	Required    bool
	Type        string
	Description string
	Value       string // Valor actual del parámetro
}

// RequestBody represents a request body
type RequestBody struct {
	Required    bool
	ContentType string
	Schema      map[string]interface{}
	Example     string
}

// APIResponse represents an API response
type APIResponse struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
	Duration   time.Duration
	Error      error
	Timestamp  time.Time
}

// OpenAPISpec represents the OpenAPI schema
type OpenAPISpec struct {
	OpenAPI    string                            `json:"openapi"`
	Info       OpenAPIInfo                       `json:"info"`
	Paths      map[string]PathItem               `json:"paths"`
	Servers    []OpenAPIServer                   `json:"servers"`
	Components map[string]map[string]interface{} `json:"components"` // For resolving $ref
}

type OpenAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type OpenAPIServer struct {
	URL         string `json:"url"`
	Description string `json:"description"`
}

type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
	Patch  *Operation `json:"patch,omitempty"`
}

type Operation struct {
	Tags        []string               `json:"tags"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Parameters  []OperationParameter   `json:"parameters"`
	RequestBody *OperationRequestBody  `json:"requestBody,omitempty"`
	Responses   map[string]interface{} `json:"responses"`
}

type OperationParameter struct {
	Name        string                 `json:"name"`
	In          string                 `json:"in"`
	Required    bool                   `json:"required"`
	Description string                 `json:"description"`
	Schema      map[string]interface{} `json:"schema"`
}

type OperationRequestBody struct {
	Required bool                       `json:"required"`
	Content  map[string]MediaTypeObject `json:"content"`
}

type MediaTypeObject struct {
	Schema  map[string]interface{} `json:"schema"`
	Example interface{}            `json:"example"`
}

// Model is the main Bubble Tea model
type Model struct {
	// Application state
	state  AppState
	width  int
	height int

	// API data
	apiURL    string
	spec      *OpenAPISpec
	endpoints []Endpoint
	baseURL   string

	// Selected endpoint
	selectedEndpoint *Endpoint
	lastResponse     *APIResponse

	// UI Components
	endpointList list.Model
	paramInputs  []textinput.Model
	currentInput int
	responseView viewport.Model
	sidebarView  viewport.Model
	listCursor   int
	listOffset   int

	// Navigation and state
	focusedPanel    int    // 0: List, 1: Editor, 2: Response
	responseContent string // pre-built highlighted content for response viewport
	clipboardMsg    string // transient "Copied!" notice shown in footer
	showHelp        bool
	err             error
	loading         bool

	// Sidebar search
	searchInput textinput.Model
	searchQuery string

	// Configuration
	styles StyleConfig

	// Endpoint grouping
	collapsedTags map[string]bool
}

// EndpointItem implements list.Item for endpoints
type EndpointItem struct {
	endpoint Endpoint
}

func (i EndpointItem) Title() string {
	return i.endpoint.Path
}

func (i EndpointItem) Description() string {
	return i.endpoint.Summary
}

func (i EndpointItem) FilterValue() string {
	return i.endpoint.Path + " " + i.endpoint.Summary
}
