<div align="center">

# ⚡ SNAG

**S**wagger **N**avigator **A**nd **G**enerator

*Cliente TUI moderno para explorar y probar APIs REST documentadas con OpenAPI 3.x / Swagger*

<br>

<table>
  <tr>
    <td><img src="images/image2.jpg" alt="SNAG Interface" width="400"/></td>
    <td><img src="images/image3.jpg" alt="SNAG Navigation" width="400"/></td>
  </tr>
  <tr>
    <td><img src="images/image4.jpg" alt="SNAG Editor" width="400"/></td>
    <td><img src="images/image5.jpg" alt="SNAG Response" width="400"/></td>
  </tr>
</table>

<br>

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![Built with Bubble Tea](https://img.shields.io/badge/Built%20with-Bubble%20Tea-FF69B4)](https://github.com/charmbracelet/bubbletea)

</div>

---

## 🌟 ¿Qué es SNAG?

SNAG es un cliente de API interactivo que se ejecuta en tu terminal. Piensa en él como Postman, pero más rápido, minimalista y con atajos de teclado estilo Vim. Consume automáticamente la especificación OpenAPI de cualquier API y te permite navegar, editar parámetros y ejecutar requests sin salir de tu terminal.

### Características Principales

- **🎯 Consumo automático de OpenAPI 3.x**: Solo proporciona la URL y SNAG descubre todos los endpoints
- **⌨️ Navegación Vim-style**: `j/k` para moverse, `enter` para seleccionar, `esc` para volver
- **🎨 Interfaz elegante tipo LazyVim**: Dark mode profesional con sintaxis highlighting
- **📋 3 Paneles interactivos en tiempo real**:
  - **Navegación**: Lista de endpoints con scroll fluido, agrupados por tags con colores por método HTTP
  - **Editor**: Sistema de campos de entrada para parámetros y request body (sin edición manual de JSON!)
  - **Respuesta**: Resultado con syntax highlighting, status code y tiempo de respuesta
- **🔍 Búsqueda/filtrado en tiempo real**: Tecla `/` para buscar endpoints por método, path o descripción
- **📊 Grupos colapsables**: Organiza endpoints por tags (Users, Products, Orders, etc.)
- **⚡ Rápido y ligero**: Sin dependencias pesadas, escrito en Go con Bubble Tea
- **📋 Copiar respuestas**: Tecla `y` para copiar el JSON al portapapeles
- **🌐 Compatible con cualquier framework**: FastAPI, Express, Spring Boot, Django, Laravel, ASP.NET Core, etc.

## 🚀 Instalación

### Opción 1: Compilar desde código fuente

```bash
# Clonar el repositorio
git clone <repo-url>
cd snag

# Instalar dependencias
go mod tidy

# Compilar
go build -o snag

# (Windows)
go build -o snag.exe
```

### Opción 2: Instalar binario (próximamente)

```bash
# Linux/macOS con Homebrew
brew install snag

# Go install
go install github.com/tu-usuario/snag@latest
```

## 📖 Uso Básico

```bash
# Con cualquier API que exponga OpenAPI 3.x
snag http://localhost:8000

# Con la ruta completa al openapi.json
snag http://localhost:8000/openapi.json

# Con API remota
snag https://api.example.com

# Con swagger.json (alias de openapi.json)
snag http://localhost:8000/swagger.json
```

### Ejemplos con diferentes frameworks

```bash
# FastAPI (Python)
snag http://localhost:8000

# Express.js con swagger-jsdoc (Node.js)
snag http://localhost:3000

# NestJS con @nestjs/swagger (Node.js)
snag http://localhost:3000

# Spring Boot con springdoc-openapi (Java)
snag http://localhost:8080

# Django con drf-spectacular (Python)
snag http://localhost:8000

# Laravel con l5-swagger (PHP)
snag http://localhost:8000

# ASP.NET Core con Swashbuckle (.NET)
snag http://localhost:5000

# Go con swaggo/swag
snag http://localhost:8080

# APIs públicas
snag https://api.stripe.com
snag https://api.github.com
```

## ⌨️ Atajos de Teclado

### Modo Navegación

- `j` / `↓` - Navegar hacia abajo
- `k` / `↑` - Navegar hacia arriba
- `enter` - Seleccionar endpoint
- `e` - Editar parámetros/body con campos de entrada
- `r` - Ejecutar request
- `tab` - Cambiar panel enfocado
- `q` - Salir

### Modo Edición de Parámetros

- `tab` / `↓` - Siguiente campo
- `shift+tab` / `↑` - Campo anterior
- `ctrl+s` - Guardar y volver
- `esc` - Cancelar y volver

### Modo Visualización de Respuesta

- `j/k` - Scroll línea a línea
- `d` - Page down
- `u` - Page up
- `g` - Ir al inicio
- `G` - Ir al final
- `r` - Re-ejecutar request
- `esc` - Volver a navegación

## 🛠️ Stack Tecnológico

- **Framework TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Arquitectura Elm para Go
- **Estilos**: [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Layouts, bordes y colores elegantes
- **Componentes**: [Bubbles](https://github.com/charmbracelet/bubbles) - Listas, textareas y viewports
- **HTTP Client**: Cliente HTTP estándar de Go
- **JSON Parsing**: Encoding/json estándar

## 📁 Arquitectura

```
snag/
├── main.go      # Entry point, Bubble Tea loop (Init, Update, View)
├── model.go     # Estructuras de datos (Model, Endpoint, OpenAPISpec, etc.)
├── api.go       # Lógica de API (fetch OpenAPI, parse endpoints, execute requests)
├── ui.go        # Estilos Lip Gloss y renderizado de paneles
└── go.mod       # Dependencias
```

### Patrón Elm (Model-View-Update)

- **Model**: Estado completo de la aplicación
- **Init**: Inicialización y carga del OpenAPI spec
- **Update**: Manejo de eventos (teclado, mensajes asíncronos)
- **View**: Renderizado de la UI con Lip Gloss

## 🎨 Tema de Colores (Dark Mode)

| Elemento            | Color                                                           | Hex       |
| ------------------- | --------------------------------------------------------------- | --------- |
| Primary (Navy Blue) | ![#1E3A8A](https://via.placeholder.com/15/1E3A8A/000000?text=+) | `#1E3A8A` |
| Secondary (Azul)    | ![#3B82F6](https://via.placeholder.com/15/3B82F6/000000?text=+) | `#3B82F6` |
| Success (Verde)     | ![#10B981](https://via.placeholder.com/15/10B981/000000?text=+) | `#10B981` |
| Error (Rojo)        | ![#EF4444](https://via.placeholder.com/15/EF4444/000000?text=+) | `#EF4444` |
| Warning (Amarillo)  | ![#F59E0B](https://via.placeholder.com/15/F59E0B/000000?text=+) | `#F59E0B` |
| Info (Cyan)         | ![#06B6D4](https://via.placeholder.com/15/06B6D4/000000?text=+) | `#06B6D4` |

### Métodos HTTP

- `GET` - Cyan (`#06B6D4`)
- `POST` - Verde (`#10B981`)
- `PUT` - Amarillo (`#F59E0B`)
- `DELETE` - Rojo (`#EF4444`)
- `PATCH` - Azul (`#3B82F6`)

## 🔮 Roadmap

- [ ] Syntax highlighting para JSON con Chroma
- [ ] Guardado de historial de requests
- [ ] Variables de entorno y templates
- [ ] Autenticación (Bearer, API Key, OAuth)
- [ ] Export de requests a curl/código
- [ ] Múltiples entornos (dev, staging, prod)
- [ ] Tests automáticos de endpoints

## 📝 Ejemplo con FastAPI

```python
# app.py
from fastapi import FastAPI

app = FastAPI(title="Mi API", version="1.0.0")

@app.get("/users")
def get_users():
    return [{"id": 1, "name": "Juan"}]

@app.post("/users")
def create_user(user: dict):
    return {"id": 2, **user}
```

```bash
# Ejecutar FastAPI
uvicorn app:app --reload

# En otra terminal, ejecutar SNAG
snag http://localhost:8000
```

## 🤝 Contribuciones

Las contribuciones son bienvenidas. Por favor, abre un issue o PR.

## 📄 Licencia

MIT

---

Hecho con ❤️ usando [Bubble Tea](https://github.com/charmbracelet/bubbletea) y [Lip Gloss](https://github.com/charmbracelet/lipgloss)
