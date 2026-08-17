package i18n

import "errors"

var (
	ErrNoLanguage    = errors.New("i18n: at least one language is required")
	ErrNoLocaleFiles = errors.New("i18n: no locale files (active.*.toml) found in the embedded filesystem")
)
