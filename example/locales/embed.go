package locales

import (
    "embed"
)

//go:embed active.*.toml
var FS embed.FS

//go:generate go run ../../cmd/i18n generate -dir . -package locales -out keys_gen.go
