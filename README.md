# kreisler/i18n

Type-safe i18n for Go, powered internally by `github.com/nicksnyder/go-i18n/v2`.

## Features

- TOML locale files.
- Generated Go constants for autocomplete.
- `go:generate` integration.
- Embedded locale files with `embed.FS`.
- Locale fallback through go-i18n.
- Named template variables.
- Plurals.
- Validation of duplicate keys and missing default-language keys.
- CLI commands: `generate` and `validate`.

## Quick start

```bash
go mod tidy
go generate ./...
go test ./...
go run ./cmd/i18n validate
```

Then:

```go
package main

import (
    "fmt"

    ki18n "github.com/itskreisler/i18n-go"
    "github.com/itskreisler/i18n-go/example/locales"
)

func main() {
    loc, err := ki18n.New(locales.FS, "en", "es")
    if err != nil {
        panic(err)
    }

    fmt.Println(loc.T(locales.Hello))
    fmt.Println(loc.TD(locales.WelcomeUser, map[string]any{
        "Name": "Kreisler",
    }))
}
```

The generated keys are in `example/locales/keys_gen.go`.

## Layout

```text
.
├── bundle.go
├── localizer.go
├── generator/
│   ├── parser.go
│   ├── generator.go
│   └── generator_test.go
├── validator/
│   └── validator.go
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

The runtime delegates message lookup, fallback, interpolation and pluralization to go-i18n. This package mainly provides a typed API and code generation around it.

`go generate` is explicit in Go; it is not run automatically by `go build` or `go test`.
