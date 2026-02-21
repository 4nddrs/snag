# ⚡ SNAG

**S**wagger **N**avigator **A**nd **G**enerator - Una herramienta TUI elegante para explorar y probar APIs documentadas con OpenAPI/Swagger.

## 🎨 Características

- **Interfaz elegante tipo LazyVim**: Dark mode con acentos azul navy
- **Navegación Vim-style**: `j/k` para moverse, `enter` para seleccionar, `esc` para volver
- **3 Paneles interactivos**:
  - 📋 **Navegación**: Lista de endpoints con scroll fluido, agrupados por tags con colores por método HTTP
  - ✏️ **Editor**: Sistema de campos de entrada para parámetros (no más edición manual de JSON!)
  - 📡 **Respuesta**: Muestra el resultado con syntax highlighting y status code
- **Feedback visual**: Bordes que cambian de color según el estado (cargando, éxito, error)
- **Consumo automático de OpenAPI**: Solo proporciona la URL y SNAG hace el resto
- **Interface en inglés**: CLI completamente en inglés para mayor accesibilidad

## 🚀 Instalación

```bash
# Clonar el repositorio
git clone <repo-url>
cd snag

# Instalar dependencias
go mod tidy

# Compilar
go build -o snag.exe
```

## 📖 Uso

```bash
# Ejecutar con la URL de tu API FastAPI
snag http://localhost:8000

# O con la ruta completa al openapi.json
snag http://localhost:8000/openapi.json
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
