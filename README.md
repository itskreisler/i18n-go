# itskreisler/i18n-go

**English** | [Español](README-es.md)

Type-safe i18n for Go, powered internally by `github.com/nicksnyder/go-i18n/v2`.

Locale files are named `active.<lang>.toml` (e.g. `active.en.toml`, `active.es.toml`)
and embedded with `//go:embed active.*.toml`. Only `active.*.toml` files are
loaded at runtime and considered by the `generate` and `validate` commands.

## Features

- TOML locale files.
- Generated Go constants for autocomplete (compile-time safety).
- `go:generate` integration.
- Embedded locale files with `embed.FS`.
- Locale fallback: a message missing from the preferred language falls back
  to the last language instead of panicking.
- Named template variables (`{{.Name}}`).
- Plurals (`zero`, `one`, `two`, `few`, `many`, `other`).
- Validation of duplicate keys and missing default-language keys.
- CLI commands: `generate` and `validate`.

## Installation

```bash
go get github.com/itskreisler/i18n-go
```

## Step-by-step guide

### 1. Create the locale files

Create a `locales/` directory in your module and add one `active.<lang>.toml`
file per language. Dotted keys become nested tables; template variables are
interpolated with `{{.Field}}`; plurals use the go-i18n categories.

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

### 2. Embed the files

`locales/embed.go`:

```go
package locales

import "embed"

//go:embed active.*.toml
var FS embed.FS

//go:generate go run github.com/itskreisler/i18n-go/cmd/i18n generate -dir . -package locales -out keys_gen.go
```

### 3. Generate the typed keys

```bash
go generate ./...
```

This creates `locales/keys_gen.go` with one typed constant per key, so typos
become compile errors instead of runtime panics:

```go
Hello         i18n.Key = "hello"
WelcomeUser   i18n.Key = "welcome_user"
Items         i18n.Key = "items"
SettingsTitle i18n.Key = "settings.title"
```

### 4. Use the localizer

```go
package main

import (
    "fmt"

    ki18n "github.com/itskreisler/i18n-go"
    "your-module/locales" // the package you created in step 2
)

func main() {
    loc, err := ki18n.New(locales.FS, "es", "en") // prefer es, fall back to en
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

### 5. Validate your translations

```bash
go run github.com/itskreisler/i18n-go/cmd/i18n validate -dir locales -default en
```

The validator reports every key missing from a locale (compared to the default)
and every extra key. It exits non-zero when something is wrong, so you can run
it in CI.

## Language fallback

`New(fsys, langs...)`:

- The **first** language is the preferred one.
- The **last** language is the bundle's default language: a message missing
  from the preferred language is resolved from the last one instead of panicking.
- A key missing from *every* loaded language panics at lookup time — with the
  generated typed keys, this should be impossible in practice.

## CLI reference

```bash
# Generate the typed keys
go run github.com/itskreisler/i18n-go/cmd/i18n generate -dir locales -package locales -out locales/keys_gen.go

# Validate that every locale defines the same keys as the default locale
go run github.com/itskreisler/i18n-go/cmd/i18n validate -dir locales -default en
```

Flags:

| Command    | Flag       | Default                     | Description               |
|------------|------------|-----------------------------|---------------------------|
| `generate` | `-dir`     | `example/locales`           | Locale directory          |
| `generate` | `-package` | `locales`                   | Go package name           |
| `generate` | `-out`     | `example/locales/keys_gen.go` | Generated output file   |
| `validate` | `-dir`     | `example/locales`           | Locale directory          |
| `validate` | `-default` | `en`                        | Default locale            |

## API

```go
func New(fsys fs.FS, langs ...string) (*Localizer, error)
func (l *Localizer) T(key Key) string
func (l *Localizer) TD(key Key, data map[string]any) string
func (l *Localizer) TP(key Key, count any, data map[string]any) string
func (l *Localizer) Localize(cfg *i18n.LocalizeConfig) (string, error)
func (l *Localizer) Bundle() *Bundle
```

Errors:

- `ErrNoLanguage` — `New` called without any language.
- `ErrNoLocaleFiles` — the filesystem contains no `active.*.toml` files.

## Example

See [`./example`](./example) for a runnable program with English and Spanish
locales (`go run ./example`).

## Layout

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

## Design

The runtime delegates message lookup, fallback, interpolation and pluralization
to go-i18n. This package mainly provides a typed API and code generation around it.

`go generate` is explicit in Go; it is not run automatically by `go build` or
`go test`.
