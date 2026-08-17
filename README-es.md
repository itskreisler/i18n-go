# itskreisler/i18n-go

i18n con seguridad de tipos para Go, impulsado internamente por
`github.com/nicksnyder/go-i18n/v2`.

Los archivos de idioma se llaman `active.<lang>.toml` (p. ej. `active.en.toml`,
`active.es.toml`) y se incrustan con `//go:embed active.*.toml`. Solo los
archivos `active.*.toml` se cargan en tiempo de ejecución y son tenidos en
cuenta por los comandos `generate` y `validate`.

## Características

- Archivos de idioma en TOML.
- Constantes Go generadas para autocompletado (seguridad en tiempo de compilación).
- Integración con `go:generate`.
- Archivos de idioma incrustados con `embed.FS`.
- Fallback de idioma: un mensaje ausente en el idioma preferido se resuelve
  desde el último idioma en lugar de provocar un panic.
- Variables de plantilla con nombre (`{{.Name}}`).
- Plurales (`zero`, `one`, `two`, `few`, `many`, `other`).
- Validación de claves duplicadas y de claves faltantes respecto al idioma
  por defecto.
- Comandos CLI: `generate` y `validate`.

## Instalación

```bash
go get github.com/itskreisler/i18n-go
```

## Guía paso a paso

### 1. Crea los archivos de idioma

Crea un directorio `locales/` en tu módulo y añade un archivo
`active.<lang>.toml` por idioma. Las claves con puntos se convierten en
tablas anidadas; las variables de plantilla se interpolan con `{{.Campo}}`;
los plurales usan las categorías de go-i18n.

`locales/active.en.toml`:

```toml
[hello]
other = "Hello!"

[welcome_user]
other = "Welcome, {{.Name}}!"

[items]
one = "You have {{.Count}} item."
other = "You have {{.Count}} items."

[settings.title]
other = "Settings"
```

`locales/active.es.toml`:

```toml
[hello]
other = "¡Hola!"

[welcome_user]
other = "¡Bienvenido, {{.Name}}!"

[items]
one = "Tienes {{.Count}} elemento."
other = "Tienes {{.Count}} elementos."

[settings.title]
other = "Configuración"
```

### 2. Incrusta los archivos

`locales/embed.go`:

```go
package locales

import "embed"

//go:embed active.*.toml
var FS embed.FS

//go:generate go run github.com/itskreisler/i18n-go/cmd/i18n generate -dir . -package locales -out keys_gen.go
```

### 3. Genera las claves tipadas

```bash
go generate ./...
```

Esto crea `locales/keys_gen.go` con una constante tipada por clave, de modo
que los errores tipográficos se convierten en errores de compilación en lugar
de panics en tiempo de ejecución:

```go
Hello         i18n.Key = "hello"
WelcomeUser   i18n.Key = "welcome_user"
Items         i18n.Key = "items"
SettingsTitle i18n.Key = "settings.title"
```

### 4. Usa el localizer

```go
package main

import (
    "fmt"

    ki18n "github.com/itskreisler/i18n-go"
    "tu-modulo/locales" // el paquete que creaste en el paso 2
)

func main() {
    loc, err := ki18n.New(locales.FS, "es", "en") // prefiere es, fallback a en
    if err != nil {
        panic(err)
    }

    fmt.Println(loc.T(locales.Hello))                             // ¡Hola!
    fmt.Println(loc.TD(locales.WelcomeUser, map[string]any{       // ¡Bienvenido, Kreisler!
        "Name": "Kreisler",
    }))
    fmt.Println(loc.TP(locales.Items, 3, map[string]any{          // Tienes 3 elementos.
        "Count": 3,
    }))
}
```

### 5. Valida tus traducciones

```bash
go run github.com/itskreisler/i18n-go/cmd/i18n validate -dir locales -default en
```

El validador informa de cada clave ausente en un idioma (en comparación con el
idioma por defecto) y de cada clave extra. Termina con código de salida
distinto de cero cuando hay problemas, por lo que puedes ejecutarlo en CI.

## Fallback de idioma

`New(fsys, langs...)`:

- El **primer** idioma es el preferido.
- El **último** idioma es el idioma por defecto del bundle: un mensaje ausente
  en el idioma preferido se resuelve desde el último en lugar de provocar un panic.
- Una clave ausente en *todos* los idiomas cargados provoca un panic al
  buscarla; con las claves tipadas generadas, esto debería ser imposible en la práctica.

## Referencia de la CLI

```bash
# Genera las claves tipadas
go run github.com/itskreisler/i18n-go/cmd/i18n generate -dir locales -package locales -out locales/keys_gen.go

# Valida que cada idioma defina las mismas claves que el idioma por defecto
go run github.com/itskreisler/i18n-go/cmd/i18n validate -dir locales -default en
```

Flags:

| Comando    | Flag       | Por defecto                   | Descripción              |
|------------|------------|-------------------------------|--------------------------|
| `generate` | `-dir`     | `example/locales`             | Directorio de idiomas    |
| `generate` | `-package` | `locales`                     | Nombre del paquete Go    |
| `generate` | `-out`     | `example/locales/keys_gen.go` | Archivo de salida generado |
| `validate` | `-dir`     | `example/locales`             | Directorio de idiomas    |
| `validate` | `-default` | `en`                          | Idioma por defecto       |

## API

```go
func New(fsys fs.FS, langs ...string) (*Localizer, error)
func (l *Localizer) T(key Key) string
func (l *Localizer) TD(key Key, data map[string]any) string
func (l *Localizer) TP(key Key, count any, data map[string]any) string
func (l *Localizer) Localize(cfg *i18n.LocalizeConfig) (string, error)
func (l *Localizer) Bundle() *Bundle
```

Errores:

- `ErrNoLanguage` — se llamó a `New` sin ningún idioma.
- `ErrNoLocaleFiles` — el sistema de archivos no contiene ningún `active.*.toml`.

## Ejemplo

Consulta [`./example`](./example) para ver un programa ejecutable con idiomas
inglés y español (`go run ./example`).

## Estructura

```text
.
├── go.mod
├── go.sum
├── bundle.go
├── errors.go
├── localizer.go
├── localizer_test.go
├── generator/
│   ├── generator.go
│   ├── generator_test.go
│   └── parser.go
├── validator/
│   ├── validator.go
│   └── validator_test.go
├── cmd/i18n/
│   └── main.go
└── example/
    ├── app.go
    └── locales/
        ├── embed.go
        ├── active.en.toml
        ├── active.es.toml
        └── keys_gen.go
```

## Diseño

El runtime delega la búsqueda de mensajes, el fallback, la interpolación y los
plurales en go-i18n. Este paquete proporciona principalmente una API tipada y
generación de código alrededor de él.

`go generate` es explícito en Go; no se ejecuta automáticamente con
`go build` ni `go test`.
