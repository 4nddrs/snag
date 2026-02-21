package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

// FetchOpenAPISpec obtiene el spec de OpenAPI desde una URL
func FetchOpenAPISpec(url string) (*OpenAPISpec, error) {
	// Asegurar que la URL termine en /openapi.json
	if !strings.HasSuffix(url, "/openapi.json") {
		url = strings.TrimSuffix(url, "/") + "/openapi.json"
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching OpenAPI spec: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var spec OpenAPISpec
	if err := json.Unmarshal(body, &spec); err != nil {
		return nil, fmt.Errorf("error parsing OpenAPI spec: %w", err)
	}

	return &spec, nil
}

// ParseEndpoints convierte el OpenAPI spec en una lista de endpoints
func ParseEndpoints(spec *OpenAPISpec) []Endpoint {
	endpoints := []Endpoint{}

	for path, pathItem := range spec.Paths {
		if pathItem.Get != nil {
			endpoints = append(endpoints, parseOperation(GET, path, pathItem.Get))
		}
		if pathItem.Post != nil {
			endpoints = append(endpoints, parseOperation(POST, path, pathItem.Post))
		}
		if pathItem.Put != nil {
			endpoints = append(endpoints, parseOperation(PUT, path, pathItem.Put))
		}
		if pathItem.Delete != nil {
			endpoints = append(endpoints, parseOperation(DELETE, path, pathItem.Delete))
		}
		if pathItem.Patch != nil {
			endpoints = append(endpoints, parseOperation(PATCH, path, pathItem.Patch))
		}
	}

	// Ordenar por tag y luego por path
	sort.Slice(endpoints, func(i, j int) bool {
		if endpoints[i].Tag != endpoints[j].Tag {
			return endpoints[i].Tag < endpoints[j].Tag
		}
		return endpoints[i].Path < endpoints[j].Path
	})

	return endpoints
}

// parseOperation convierte una Operation en un Endpoint
func parseOperation(method HTTPMethod, path string, op *Operation) Endpoint {
	endpoint := Endpoint{
		Method:      method,
		Path:        path,
		Summary:     op.Summary,
		Description: op.Description,
		Parameters:  []Parameter{},
	}

	// Determinar tag (usar el primero si existe, si no usar "default")
	if len(op.Tags) > 0 {
		endpoint.Tag = op.Tags[0]
	} else {
		endpoint.Tag = "default"
	}

	// Parsear parámetros
	for _, param := range op.Parameters {
		paramType := "string" // tipo por defecto
		if schema, ok := param.Schema["type"].(string); ok {
			paramType = schema
		}

		endpoint.Parameters = append(endpoint.Parameters, Parameter{
			Name:        param.Name,
			In:          param.In,
			Required:    param.Required,
			Type:        paramType,
			Description: param.Description,
		})
	}

	// Parsear request body
	if op.RequestBody != nil {
		for contentType, mediaType := range op.RequestBody.Content {
			example := "{}"
			if mediaType.Example != nil {
				exampleJSON, _ := json.MarshalIndent(mediaType.Example, "", "  ")
				example = string(exampleJSON)
			} else if mediaType.Schema != nil {
				// Generar ejemplo básico desde schema
				example = generateExampleFromSchema(mediaType.Schema)
			}

			endpoint.RequestBody = &RequestBody{
				Required:    op.RequestBody.Required,
				ContentType: contentType,
				Schema:      mediaType.Schema,
				Example:     example,
			}
			break // Solo tomar el primer content type
		}
	}

	return endpoint
}

// generateExampleFromSchema genera un JSON de ejemplo básico desde un schema
func generateExampleFromSchema(schema map[string]interface{}) string {
	// Implementación básica - puede ser expandida
	if schema == nil {
		return "{}"
	}

	schemaType, ok := schema["type"].(string)
	if !ok {
		return "{}"
	}

	switch schemaType {
	case "object":
		properties, ok := schema["properties"].(map[string]interface{})
		if !ok {
			return "{}"
		}

		example := map[string]interface{}{}
		for key, prop := range properties {
			propMap, ok := prop.(map[string]interface{})
			if !ok {
				continue
			}
			propType, _ := propMap["type"].(string)
			switch propType {
			case "string":
				example[key] = "string"
			case "integer":
				example[key] = 0
			case "number":
				example[key] = 0.0
			case "boolean":
				example[key] = false
			case "array":
				example[key] = []interface{}{}
			case "object":
				example[key] = map[string]interface{}{}
			}
		}

		exampleJSON, _ := json.MarshalIndent(example, "", "  ")
		return string(exampleJSON)
	}

	return "{}"
}

// ExecuteRequest ejecuta una petición HTTP a un endpoint
func ExecuteRequest(baseURL string, endpoint *Endpoint, body string) (*APIResponse, error) {
	startTime := time.Now()

	// Construir URL completa
	url := baseURL + endpoint.Path

	// Reemplazar path parameters
	for _, param := range endpoint.Parameters {
		if param.In == "path" && param.Value != "" {
			url = strings.Replace(url, "{"+param.Name+"}", param.Value, 1)
		}
	}

	// Añadir query parameters
	queryParams := []string{}
	for _, param := range endpoint.Parameters {
		if param.In == "query" && param.Value != "" {
			queryParams = append(queryParams, fmt.Sprintf("%s=%s", param.Name, param.Value))
		}
	}
	if len(queryParams) > 0 {
		url += "?" + strings.Join(queryParams, "&")
	}

	// Crear request
	var req *http.Request
	var err error

	if body != "" && endpoint.RequestBody != nil {
		req, err = http.NewRequest(string(endpoint.Method), url, bytes.NewBufferString(body))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
		req.Header.Set("Content-Type", endpoint.RequestBody.ContentType)
	} else {
		req, err = http.NewRequest(string(endpoint.Method), url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
	}

	// Añadir header parameters
	for _, param := range endpoint.Parameters {
		if param.In == "header" && param.Value != "" {
			req.Header.Set(param.Name, param.Value)
		}
	}

	// Ejecutar request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return &APIResponse{
			Error:     err,
			Duration:  time.Since(startTime),
			Timestamp: time.Now(),
		}, err
	}
	defer resp.Body.Close()

	// Leer respuesta
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &APIResponse{
			StatusCode: resp.StatusCode,
			Error:      err,
			Duration:   time.Since(startTime),
			Timestamp:  time.Now(),
		}, err
	}

	// Formatear JSON si es posible
	formattedBody := string(respBody)
	var jsonData interface{}
	if json.Unmarshal(respBody, &jsonData) == nil {
		formatted, err := json.MarshalIndent(jsonData, "", "  ")
		if err == nil {
			formattedBody = string(formatted)
		}
	}

	return &APIResponse{
		StatusCode: resp.StatusCode,
		Body:       formattedBody,
		Headers:    resp.Header,
		Duration:   time.Since(startTime),
		Timestamp:  time.Now(),
	}, nil
}

// GetBaseURL obtiene la URL base del primer servidor en el spec
func GetBaseURL(spec *OpenAPISpec, apiURL string) string {
	if len(spec.Servers) > 0 {
		return spec.Servers[0].URL
	}

	// Si no hay servers definidos, usar la URL base del API
	if strings.HasSuffix(apiURL, "/openapi.json") {
		return strings.TrimSuffix(apiURL, "/openapi.json")
	}

	return apiURL
}
